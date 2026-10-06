# Fonts

The files here are test fixtures: no part of AntUI loads them at runtime, and
a stylesheet only reaches one of them when a test writes an @font-face for it.

## Roboto-Regular.ttf, Roboto-Bold.ttf, Roboto-Italic.ttf

The Roboto family, Copyright 2015 Google Inc., taken from the Debian
`fonts-roboto-unhinted` package (2:0~20170802-3), which carries the unhinted
TrueType builds of the family.

Roboto is licensed under the Apache License, Version 2.0 — the same license
this repository is under; see `LICENSE` at the root of it, or
https://www.apache.org/licenses/LICENSE-2.0

SPDX identifier: Apache-2.0
Upstream: https://github.com/googlefonts/roboto-classic

## DejaVuSansMono.ttf

DejaVu Sans Mono, from the DejaVu font family, taken from the Debian
`fonts-dejavu-core` package. DejaVu's changes are in the public domain; the
typeface itself carries the Bitstream Vera license, whose permission notice
must travel with copies of it and is reproduced below. DejaVu is the
second family a test's font-family list falls back through — its glyph set
differs from Roboto's, which is what a fallback test needs.

Source: https://dejavu-fonts.github.io/

Copyright (c) 2003 by Bitstream, Inc. All Rights Reserved.
Bitstream Vera is a trademark of Bitstream, Inc.
DejaVu changes are in the public domain.

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of the fonts accompanying this license ("Fonts") and associated
> documentation files (the "Font Software"), to reproduce and distribute the
> Font Software, including without limitation the rights to use, copy, merge,
> publish, distribute, and/or sell copies of the Font Software, and to permit
> persons to whom the Font Software is furnished to do so, subject to the
> following conditions:
>
> The above copyright and trademark notices and this permission notice shall
> be included in all copies of one or more of the Font Software typefaces.
>
> The Font Software may be modified, altered, or added to, and in particular
> the designs of glyphs or characters in the Fonts may be modified and
> additional glyphs or characters may be added to the Fonts, only if the fonts
> are renamed to names not containing either the words "Bitstream" or the word
> "Vera".
>
> This License becomes null and void to the extent applicable to Fonts or
> Font Software that has been modified and is distributed under the "Bitstream
> Vera" names.
>
> The Font Software may be sold as part of a larger software package but no
> copy of one or more of the Font Software typefaces may be sold by itself.
>
> THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
> OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT OF COPYRIGHT, PATENT,
> TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL BITSTREAM OR THE GNOME
> FOUNDATION BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, INCLUDING
> ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL DAMAGES,
> WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF
> THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM OTHER DEALINGS IN THE
> FONT SOFTWARE.
>
> Except as contained in this notice, the names of Gnome, the Gnome
> Foundation, and Bitstream Inc., shall not be used in advertising or
> otherwise to promote the sale, use or other dealings in this Font Software
> without prior written authorization from the Gnome Foundation or Bitstream
> Inc. For further information, contact: fonts at gnome dot org.

SPDX identifier: Bitstream-Vera
