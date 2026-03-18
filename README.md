# Compile for alpine

```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
```

# Upload resources
```shell
scp main diabasis:main2
tar czf assets.tgz assets
scp assets.tgz root@10.0.10.19:assets.tgz
ssh root@10.0.10.19 tar -xzf assets.tgz 
ssh root@10.0.10.19 rm assets.tgz
ssh root@10.0.10.19 mv main2 main
ssh root@10.0.10.19 pkill main
rm assets.tgz
````


# Setup to run in Alpine Linux
Edit GRUB, update and reboot
```shell
vi /etc/default/grub

GRUB_TIMEOUT=2
GRUB_DISABLE_SUBMENU=y
GRUB_DISABLE_RECOVERY=true
GRUB_CMDLINE_LINUX_DEFAULT="modules=sd-mod,usb-storage,ext4 quiet rootfstype=ext4 video=efifb:320x200 nomodeset"
GRUB_GFXMODE=320x200
GRUB_GFXPAYLOAD_LINUX=keep

grub-mkconfig -o /boot/grub/grub.cfg
reboot
```


edit /etc/inittab and add the following line to run the program on tty1

```shell
tty1::respawn:/bin/sh -c "cd /root && exec ./main"
```
