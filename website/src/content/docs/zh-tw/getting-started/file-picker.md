---
title: 檔案選擇器
description: 透過 --chooser-file 讓 superfile 為腳本、編輯器和桌面應用程式選擇檔案。
head:
  - tag: title
    content: 檔案選擇器 | superfile
---

superfile 可以當作檔案選擇器使用。以 `--chooser-file <路徑>`（簡寫：`-cf`）啟動後，superfile 不會開啟你選擇的項目，而是把選取的路徑寫入該檔案並結束。

```bash
spf --chooser-file /tmp/picked ~/Downloads
cat /tmp/picked
```

## 寫入的內容

| 操作                                                    | 寫入 chooser 檔案的內容        |
| ------------------------------------------------------- | ------------------------------ |
| 在檔案上按 `enter`，或按 `e`（`open_file_with_editor`） | 目前的檔案，或所有已選取的項目 |
| `E`（`open_current_directory_with_editor`）             | 目前所在的目錄                 |

多個路徑會一行一個寫入。使用選取模式（`v`）可以選擇多個項目。若 chooser 檔案無法寫入，superfile 會記錄錯誤並照常開啟該項目。

## 在 shell 腳本中使用

```bash
picked="$(mktemp)"
spf --chooser-file "$picked"
while IFS= read -r file; do
    echo "You picked: $file"
done < "$picked"
rm -f "$picked"
```

## 在 Linux 上作為桌面檔案選擇器

[xdg-desktop-portal-termfilechooser](https://github.com/hunkyburrito/xdg-desktop-portal-termfilechooser) 讓透過桌面 portal 選擇檔案的應用程式（Firefox、Chromium、Discord 等）改用終端機檔案管理器。它附帶一個執行 `spf --chooser-file` 的 `superfile-wrapper.sh`。

安裝 portal 後，在 `~/.config/xdg-desktop-portal-termfilechooser/config` 中設定包裝腳本與終端機：

```ini
[filechooser]
cmd=superfile-wrapper.sh
default_dir=$HOME
env=TERMCMD=foot -T "terminal filechooser"
```

如何在 `portals.conf` 中啟用它以及其他終端機的設定，請參考該 portal 的 README。
