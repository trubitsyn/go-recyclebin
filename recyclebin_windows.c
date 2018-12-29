// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

#include <windows.h>
#include <Shobjidl.h>
#include <VersionHelpers.h>

#include "recyclebin_windows.h"

void move_to_trash(const char *filename)
{
    if (IsWindowsVistaOrGreater())
    {
        // IFileOperation
    }
    else
    {
        SHFILEOPSTRUCT operation;
    	SHFileOperation(&operation);
    }
}

void restore_from_trash(const char *filename)
{
    if (IsWindowsVistaOrGreater())
    {
        // IFileOperation
    }
    else
    {
        // SHFILEOPSTRUCT
    }
}

void delete_from_trash(const char *filename)
{
    if (IsWindowsVistaOrGreater())
    {
        // IFileOperation
    }
    else
    {
        // SHFILEOPSTRUCT
    }
}

void empty_trash(void)
{
    SHEmptyRecycleBinW(NULL, NULL, SHERB_NOCONFIRMATION | SHERB_NOPROGRESSUI | SHERB_NOSOUND);
}
