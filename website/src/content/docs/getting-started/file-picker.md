---
title: File Picker
description: Use superfile to pick files for scripts, editors and desktop apps with --chooser-file.
head:
  - tag: title
    content: File Picker | superfile
---

superfile can act as a file picker. Start it with `--chooser-file <path>` (short: `-cf`). Instead of opening what you pick, superfile writes the chosen paths to that file and exits.

```bash
spf --chooser-file /tmp/picked ~/Downloads
cat /tmp/picked
```

## What gets written

| Action                                              | Written to the chooser file              |
| --------------------------------------------------- | ---------------------------------------- |
| `enter` on a file, or `e` (`open_file_with_editor`) | The focused file, or every selected item |
| `E` (`open_current_directory_with_editor`)          | The current directory                    |

Multiple paths are written one per line. Use select mode (`v`) to pick more than one item. If the chooser file can't be written, superfile logs the error and opens the item as usual.

## Use it in a shell script

```bash
picked="$(mktemp)"
spf --chooser-file "$picked"
while IFS= read -r file; do
    echo "You picked: $file"
done < "$picked"
rm -f "$picked"
```

## Use it as the desktop file picker on Linux

[xdg-desktop-portal-termfilechooser](https://github.com/hunkyburrito/xdg-desktop-portal-termfilechooser) lets apps that use the desktop portal (Firefox, Chromium, Discord, ...) open a terminal file manager when they ask for a file. It ships a `superfile-wrapper.sh` that runs `spf --chooser-file`.

After installing the portal, set the wrapper and your terminal in `~/.config/xdg-desktop-portal-termfilechooser/config`:

```ini
[filechooser]
cmd=superfile-wrapper.sh
default_dir=$HOME
env=TERMCMD=foot -T "terminal filechooser"
```

See the portal's README for enabling it in `portals.conf` and for other terminals.
