#!/bin/bash

Xvfb :99 -screen 0 640x480x8 -nolisten tcp & google-chrome & /app/service