/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package unix

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrMissingDirEntry = errors.New("missing dir entry, probably socket files not found")
)

type socketDialler struct {
	l *slog.Logger
	e errorFormatterService

	dirName     string
	filePattern string

	dirEntries           []os.DirEntry
	count                uint
	currentEntryPosition uint
}

func (d *socketDialler) next() (os.DirEntry, bool) {
	position := d.currentEntryPosition
	d.currentEntryPosition++

	if d.currentEntryPosition <= d.count {
		err := d.prepare()
		if err != nil {
			return nil, false
		}

		return d.dirEntries[position], false
	}

	return d.dirEntries[position], true
}

func (d *socketDialler) Prepare() (func(context.Context, string) (net.Conn, error), error) {
	err := d.prepare()
	if err != nil {
		return nil, err
	}

	return d.DialCallback, nil
}

func (d *socketDialler) prepare() error {
	for {
		count, err := d.reset()
		if err == nil && count > 0 {
			return nil
		}

		time.Sleep(time.Second)
	}
}

func (d *socketDialler) reset() (uint, error) {
	d.dirEntries = make([]os.DirEntry, 0)
	d.count = 0
	d.currentEntryPosition = 0

	files, err := os.ReadDir(d.dirName)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		match, loopErr := filepath.Match(d.filePattern, file.Name())
		if loopErr != nil {
			return 0, loopErr
		}

		if match {
			d.dirEntries = append(d.dirEntries, file)
			d.count++
		}
	}

	return d.count, nil
}

func (d *socketDialler) DialCallback(ctx context.Context, _ string) (net.Conn, error) {
	var connected = false
	var conn *net.UnixConn = nil

	file, hasNext := d.next()
	hasNext = true //hack for first loop iteration

	for file != nil && !hasNext {
		filePath := filepath.Join(d.dirName, file.Name())

		resolved, err := net.ResolveUnixAddr("unix", filePath)
		if err != nil {
			d.l.Error("unable to resolve unix-socket",
				slog.String("error", err.Error()))

			file, hasNext = d.next()

			continue
		}

		dialConn, err := net.DialUnix("unix", nil, resolved)
		if err != nil {
			d.l.Error("unable to dial unix-socket",
				slog.String("error", err.Error()))

			file, hasNext = d.next()

			continue
		}

		connected = true
		conn = dialConn

		break
	}

	if conn == nil || !connected {
		return nil, ErrMissingDirEntry
	}

	return conn, nil
}

func NewUnitFileSocketDialer(logger *slog.Logger,
	errFmtSvc errorFormatterService,
	dirName,
	filePattern string,
) *socketDialler {
	return &socketDialler{
		l:           logger,
		e:           errFmtSvc,
		dirName:     dirName,
		filePattern: filePattern,

		dirEntries:           nil,
		count:                0,
		currentEntryPosition: 0,
	}
}
