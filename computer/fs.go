package computer

import "context"

type FileSystem interface {
	//IsDriveRoot returns true if a path is mounted to the parent filesystem. The root filesystem "/" is
	//considered a mount, along with disk folders and the rom folder. Other programs (such as network shares)
	//can exstend this to make other mount types by correctly assigning their return value for getDrive.
	IsDriveRoot(ctx context.Context, path string) (bool, error)

	//Complete provides completion for a file or directory name, suitable for use with _G.read.
	Complete(ctx context.Context, path, location string, includeFiles, includeDirs bool) ([]string, error)

	//List returns a list of files in a directory.
	List(ctx context.Context, path string) ([]string, error)
}
