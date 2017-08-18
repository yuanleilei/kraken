package torlib

import (
	"io/ioutil"
	"math/rand"
	"os"

	"code.uber.internal/infra/kraken/utils/randutil"
)

const fixtureTempDir = "/tmp/kraken_fixtures"

func init() {
	os.Mkdir(fixtureTempDir, 0755)
}

// PeerIDFixture returns a randomly generated PeerID.
func PeerIDFixture() PeerID {
	var p PeerID
	if _, err := rand.Read(p[:]); err != nil {
		panic(err)
	}
	return p
}

// InfoHashFixture returns a randomly generated InfoHash.
func InfoHashFixture() InfoHash {
	return MetaInfoFixture().InfoHash
}

// PeerInfoFixture returns a randomly generated PeerInfo.
func PeerInfoFixture() *PeerInfo {
	return &PeerInfo{
		InfoHash:        InfoHashFixture().String(),
		PeerID:          PeerIDFixture().String(),
		IP:              randutil.IP(),
		Port:            int64(randutil.Port()),
		DC:              "sjc1",
		BytesDownloaded: rand.Int63n(4096),
		BytesUploaded:   rand.Int63n(4096),
		BytesLeft:       rand.Int63n(4096),
	}
}

// PeerInfoForMetaInfoFixture returns a randomly generated PeerInfo associated
// with the given MetaInfo.
func PeerInfoForMetaInfoFixture(mi *MetaInfo) *PeerInfo {
	p := PeerInfoFixture()
	p.InfoHash = mi.GetInfoHash().String()
	p.BytesDownloaded = rand.Int63n(mi.Info.Length)
	p.BytesUploaded = rand.Int63n(mi.Info.Length)
	p.BytesLeft = rand.Int63n(mi.Info.Length)
	return p
}

// TestTorrentFile joins a MetaInfo with the file contents used to generate
// said MetaInfo. Note, does not include any physical files so no cleanup is
// necessary.
type TestTorrentFile struct {
	MetaInfo *MetaInfo
	Content  []byte
}

// CustomTestTorrentFileFixture returns a randomly generated TestTorrentFile
// of the given size and piece length.
// TODO(codyg): Move this to storage package.
func CustomTestTorrentFileFixture(size int, pieceLength int) *TestTorrentFile {
	f, err := ioutil.TempFile(fixtureTempDir, "torrent_")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	content := randutil.Text(size)
	if err := ioutil.WriteFile(f.Name(), content, 0755); err != nil {
		panic(err)
	}
	info, err := NewInfoFromFile(f.Name(), int64(pieceLength))
	if err != nil {
		panic(err)
	}
	mi, err := NewMetaInfoFromInfo(info, "")
	if err != nil {
		panic(err)
	}
	return &TestTorrentFile{mi, content}
}

// TestTorrentFileFixture returns a randomly generated TestTorrentFile.
func TestTorrentFileFixture() *TestTorrentFile {
	return CustomTestTorrentFileFixture(128, 32)
}

// MetaInfoFixture returns a randomly generated MetaInfo.
func MetaInfoFixture() *MetaInfo {
	return TestTorrentFileFixture().MetaInfo
}