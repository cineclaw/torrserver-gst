package torrstor

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anacrolix/torrent/metainfo"
	torrenttypes "github.com/anacrolix/torrent/types"

	"server/settings"
)

func TestSelectiveIndexPieceProtectionAndPreload(t *testing.T) {
	tempDir := t.TempDir()
	oldPath := settings.Path
	settings.Path = tempDir
	defer func() { settings.Path = oldPath }()

	const pieceLen = int64(1024 * 1024) // 1MB pieces
	const pieceCount = 30               // 30MB total torrent

	// Mock 2 files:
	// File 1: 15MB (pieces 0..14) -> head 8MB = 0..8, tail 8MB = 7..14
	// File 2: 15MB (pieces 15..29) -> head 8MB = 15..23, tail 8MB = 22..29
	info := &metainfo.Info{
		Name:        "TestSeason",
		PieceLength: pieceLen,
		Pieces:      make([]byte, pieceCount*20),
		Files: []metainfo.FileInfo{
			{Length: 15 * pieceLen, Path: []string{"S01E01.mkv"}},
			{Length: 15 * pieceLen, Path: []string{"S01E02.mkv"}},
		},
	}
	hash := metainfo.Hash([20]byte{1, 2, 3, 4, 5})

	c := &Cache{
		capacity:      10 * pieceLen, // Small capacity (10MB) to force eviction
		pieces:        make(map[int]*Piece),
		readers:       make(map[*Reader]struct{}),
		cleanTrigger:  make(chan struct{}, 1),
		cleanStop:     make(chan struct{}),
		localPriority: make(map[int]torrenttypes.PiecePriority),
	}
	c.Init(info, hash)
	defer close(c.cleanStop)

	// Verify that index pieces are marked
	if !c.isIndexPiece[0] {
		t.Fatalf("piece 0 should be index piece")
	}
	if !c.isIndexPiece[15] {
		t.Fatalf("piece 15 (start of Ep 2) should be index piece")
	}
	if !c.isIndexPiece[29] {
		t.Fatalf("piece 29 (end of Ep 2) should be index piece")
	}

	// Verify protection in isIdInFileBE even when no active reader covers Ep 2
	ranges := []Range{} // no active readers
	if !c.isIdInFileBE(ranges, 0) {
		t.Errorf("piece 0 should be protected by isIdInFileBE")
	}
	if !c.isIdInFileBE(ranges, 15) {
		t.Errorf("piece 15 should be protected by isIdInFileBE")
	}

	// Test writing and disk caching for an index piece
	p0 := c.pieces[0]
	testPayload := make([]byte, pieceLen)
	for i := range testPayload {
		testPayload[i] = 0x42
	}
	_, err := p0.WriteAt(testPayload, 0)
	if err != nil {
		t.Fatalf("WriteAt failed: %v", err)
	}
	_ = p0.MarkComplete()

	// Wait briefly for background disk write
	cacheFile := filepath.Join(tempDir, "index_cache", hash.HexString(), "piece_0.bin")
	var fileData []byte
	for i := 0; i < 20; i++ {
		if data, err := os.ReadFile(cacheFile); err == nil && len(data) == int(pieceLen) {
			fileData = data
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if len(fileData) != int(pieceLen) {
		// manually invoke if background write didn't finish
		p0.saveIndexPieceToDisk()
		fileData, err = os.ReadFile(cacheFile)
		if err != nil {
			t.Fatalf("Failed to save index piece to disk: %v", err)
		}
	}
	if fileData[0] != 0x42 {
		t.Errorf("Disk cache content mismatch")
	}

	// Now test Preload in a fresh Cache instance
	c2 := &Cache{
		capacity:      10 * pieceLen,
		pieces:        make(map[int]*Piece),
		readers:       make(map[*Reader]struct{}),
		cleanTrigger:  make(chan struct{}, 1),
		cleanStop:     make(chan struct{}),
		localPriority: make(map[int]torrenttypes.PiecePriority),
	}
	c2.Init(info, hash)
	defer close(c2.cleanStop)

	p0_loaded := c2.pieces[0]
	if !p0_loaded.Complete.Load() {
		t.Errorf("piece 0 should be marked complete after preload")
	}
	if atomic.LoadInt64(&p0_loaded.Size) != pieceLen {
		t.Errorf("piece 0 size should be %d, got %d", pieceLen, atomic.LoadInt64(&p0_loaded.Size))
	}
	readBuf := make([]byte, 10)
	n, _ := p0_loaded.ReadAt(readBuf, 0)
	if n != 10 || readBuf[0] != 0x42 {
		t.Errorf("ReadAt after preload returned incorrect data: %v", readBuf)
	}
}
