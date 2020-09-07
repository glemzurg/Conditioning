# Conditioning
Slide show software for conditioning affirmations.

Commands are:

* LEFT+RIGHT ARROWS (go back and forth in slide show)
* SPACE (start/stop slide show)
* R (toggle order of slideshow to random)
* L (reload from file)

Speed of the slide show is set in the config.json.

# Building

This software is written in golang 1.14 and built with the GTK.

Golang is pretty simple to install (https://golang.org/doc/install), GTK a little harder, follow the steps here: https://github.com/gotk3/gotk3

The crux of it is that some package need to exist in ubuntu for the build to work, and the very first time it builds it will take quite a while.

Critically there are two values that need to be set in the local environment:

* GOPATH is where golang will put temporary files that it works with.
* GOBIN (usually $GOPATH/bin) is where the binary executable will be built.

The ruh.sh scripts in this repo rely on GOBIN being set in order to find the executable.

The root-level run.sh will do a build then run so, once the necessary parts are on the local machine the whole software can be launched from

    ./run.sh

(If you cannot execute the shell script, give it executable privileges: sudo chmod a+x run.sh)

# The Affirmations File.

Every line in the affirmatinos file is either:

* An affirmation.
* A comment that is ignored (line begins with "//").
* A blank line.

For affirmations, the first part of the line is the affirmation.

The second part of the line, between "[" and "]" indicates the font, color, and positioning of the affirmation; and the image, image size, and image position.

Those parts are optional. If both the font is configured and the image is configured, their relative settings are separated by a space.

The font settings are:

* color, "b" (black with white outline) or "w" (white with black outline).
* font size, a number.
* offset from center, an x, y value with negative going up and to the left, and positive going down and to the right.

The images settings are:

* image name without path, expects the file to be in the images folder next to the affirmations file used
* image size, a decimal value where 1.0 is 100% size of the image
* offset from center, an x, y value with negative going up and to the left, and positive going down and to the right.

The affirmations.example.txt shows examples of all these setttings.

# The Configurations File.

The Configurations file is a JSON formatted file that has settings like the screen title, slide show speed, and default fonts and outline widths.

If you're not familiar with JSON files they are meant to be modified by hand but can be finicky. It is sometimes handy to put a misbehaving file in an online JSON parser and let the parser point out the place in the file that is barfing.

# The examples run.sh

The run.sh in examples show a setup that allows multiple slide shows using a single config. The run.sh expects to find the executable in the $GOBIN path and takes a subdirectory (a specific slide show) as a parameter. It is executed like this:

    cd examples
    ./run.sh slideshow1

# Royalty Free Stock Images

Thanks to https://www.pexels.com/ and the artist there for their royalty free stock images in the examples.
