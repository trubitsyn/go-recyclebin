// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

#include <windows.h>
#include "recyclebin_windows.h"

void move_to_trash(char *filename)
{
	SHFILEOPSTRUCT operation;

	SHFileOperation(&operation);
}

void restore_from_trash(char *filename)
{
}

void delete_from_trash(char *filename)
{
}

void empty_trash()
{
	SHEmptyRecycleBin(NULL, NULL, SHERB_NOCONFIRMATION | SHERB_NOPROGRESSUI | SHERB_NOSOUND);
}
