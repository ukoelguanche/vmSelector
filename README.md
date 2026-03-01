# Compile for alpine

```shell
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o fbtest .
$ scp fbtest root@10.0.10.19:fbtest
```

# Upload resources
```shell
scp resources/sprites/hud.png root@10.0.10.19:resources/sprites/HUD.png
scp resources/sprites/hud.json root@10.0.10.19:resources/sprites/HUD.json

````




# Alpine setup
Edit GRUB, update an reboot
```shell
$ vi /etc/default/grub

GRUB_TIMEOUT=2
GRUB_DISABLE_SUBMENU=y
GRUB_DISABLE_RECOVERY=true
GRUB_CMDLINE_LINUX_DEFAULT="modules=sd-mod,usb-storage,ext4 quiet rootfstype=ext4 video=efifb:320x200 nomodeset"
GRUB_GFXMODE=320x200
GRUB_GFXPAYLOAD_LINUX=keep

$ grub-mkconfig -o /boot/grub/grub.cfg
$ reboot
```

Enable fb output
```shell
$ echo 0 > /sys/class/vtconsole/vtcon1/bind
```

Test with noise 
```shell
$ cat /dev/urandom > /dev/fb0
```
