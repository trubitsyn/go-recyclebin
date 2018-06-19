// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

#ifndef RECYCLEBIN_H
#define RECYCLEBIN_H

void move_to_trash(char *filename);

void restore_from_trash(char *filename);

void delete_from_trash(char *filename);

void empty_trash();

#endif
