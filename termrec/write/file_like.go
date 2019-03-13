package write

// FileLike is a small interface to act like an os.File, for the purposes of StreamingFileWriter,
// and in particular for unit testing the StreamingFileWriter.
type FileLike interface {
	Close() error
	Name() string
}
