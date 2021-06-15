#!/bin/sh
# This is for cleaning your setup (For upgrades maybe)
# If setted up using defaults

echo Removing binaries
# Remove binary
rm -f /usr/bin/diploy

echo Removing logs and database
# Clean logs and database
rm -rf /var/log/diploy

echo Removing systemd entry
# Remove systemd entry
systemctl stop diploy
systemctl disable diploy
rm -f /etc/systemd/system/diploy.service