/*
(BSD 2-clause license)

Copyright (c) 2014, Shawn Webb
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

   * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package zfs

import (
    "os/exec"
    "strings"
)

type Dataset struct {
    DatasetPath string
    Snapshots []string
    Options map[string]string
    Mountpoint string
    Children []*Dataset
}

func byteToString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }

        n = i
    }

    return string(c[:n+1])
}

func GetDataset(name string) *Dataset {
    cmd := exec.Command("zfs", "list", "-r", "-t", "all", "-o", "name", "-d", "1", name)
    rawoutput, err := cmd.Output()

    if err != nil {
        panic(err)
        return nil
    }

    dataset := new(Dataset)
    lines := strings.Split(byteToString(rawoutput), "\n")
    for i := range lines {
        if lines[i] == name {
            continue
        }

        if strings.Contains(lines[i], "@") == false {
            dataset.Children = append(dataset.Children, GetDataset(strings.Trim(lines[i], "\n")))
            continue
        }

        dataset.Snapshots = append(dataset.Snapshots, strings.Trim(lines[i], "\n"))
    }

    return dataset
}

func GetDatasetPath(dataset string) (string, error) {
    cmd := exec.Command("/sbin/zfs", "get", "-H", "-o", "value", "mountpoint", dataset)
    output, err := cmd.CombinedOutput()

    if (err != nil) {
        return "", err
    }

    return byteToString(output), nil
}
