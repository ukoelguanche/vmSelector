# Compile for alpine

```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
```

# Upload resources
```shell
scp fbtest root@10.0.10.19:fbtest
tar czf resources.tgz resources
scp resources.tgz root@10.0.10.19:resources.tgz
ssh  root@10.0.10.19 tar -xzf resources.tgz 
````


# Alpine setup
Edit GRUB, update an reboot
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

Enable fb output
```shell
echo 0 > /sys/class/vtconsole/vtcon1/bind
```

Test with noise 
```shell
cat /dev/urandom > /dev/fb0
```


Skip library intermediates:
GOPRIVATE=github.com/ukoelguanche/graphicsengine go get github.com/ukoelguanche/graphicsengine@v0.1.3
