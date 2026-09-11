# uconsole-installer

>Fuzzy search and install for the 13 000+ Arch Linux ARM packages available to the uConsole.

Native GTK3 frontend over pacman. Loads every package from the synced `pacman -Ss` index in under a second, ranks matches with exact > prefix > substring on the name, and falls back to the description. Enter on a row pops a floating terminal with the `pacman -S` or `pacman -R` command.

## Demo

```
[screenshot placeholder]
```

## Dependencies

- GTK 3.24, pacman, Sway/Wayland
- `bash` o `sh` para el `install.sh`
- glibc estándar

## Install

Una línea, vía curl al instalador del repo:

```sh
curl -sSL https://raw.githubusercontent.com/4rgs-org/uconsole-installer/main/install.sh | sh
```


Para una versión específica:

```sh
curl -sSL https://raw.githubusercontent.com/4rgs-org/uconsole-installer/main/install.sh | sh -s -- --tag v0.1.0
```

Para instalar en otro directorio (útil para `~/.local/bin` sin root):

```sh
INSTALL_DIR=$HOME/.local/bin REPO=uconsole-installer bash install.sh
```

El instalador:
- Detecta la arquitectura con `uname -m` (aarch64, x86_64, armv7).
- Baja el tarball de la release seleccionada.
- Verifica `SHA256SUMS` antes de extraer.
- Copia el binario a `/usr/local/bin/uconsole-installer` (o el destino elegido).

## Post-install

`uconsole-installer` se lanza desde *uConsole → Search packages*.
Si lo querés invocar de forma independiente:

```sh
/usr/local/bin/uconsole-installer
```

## Build from source

```sh
git clone https://github.com/4rgs-org/uconsole-installer
cd uconsole-installer
go build -o uconsole-installer .
sudo install -m 0755 uconsole-installer /usr/local/bin/uconsole-installer
```

## Usage

| Acción | Resultado |
| --- | --- |
| `Type any substring` | Filtra en memoria; ranking con preferencia por nombre exacto. |
| ``Enter`` | Lanza `pacman -S <paquete>` o `pacman -R <paquete>` en popup flotante. |
| ``Esc`` | Cierra la ventana. |

## License

MIT — ver [LICENSE](LICENSE).

## Liability

Software is provided as-is. The maintainers are not responsible for loss of work or data caused by accidental command execution. Always confirm the menu entry before pressing Enter.
