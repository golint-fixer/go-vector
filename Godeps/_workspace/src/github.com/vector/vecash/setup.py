#!/usr/bin/env python
import os
from distutils.core import setup, Extension
sources = [
    'src/python/core.c',
    'src/libvecash/io.c',
    'src/libvecash/internal.c',
    'src/libvecash/sha3.c']
if os.name == 'nt':
    sources += [
        'src/libvecash/util_win32.c',
        'src/libvecash/io_win32.c',
        'src/libvecash/mmap_win32.c',
    ]
else:
    sources += [
        'src/libvecash/io_posix.c'
    ]
depends = [
    'src/libvecash/vecash.h',
    'src/libvecash/compiler.h',
    'src/libvecash/data_sizes.h',
    'src/libvecash/endian.h',
    'src/libvecash/vecash.h',
    'src/libvecash/io.h',
    'src/libvecash/fnv.h',
    'src/libvecash/internal.h',
    'src/libvecash/sha3.h',
    'src/libvecash/util.h',
]
pyvecash = Extension('pyvecash',
                     sources=sources,
                     depends=depends,
                     extra_compile_args=["-Isrc/", "-std=gnu99", "-Wall"])

setup(
    name='pyvecash',
    author="Matthew Wampler-Doty",
    author_email="matthew.wampler.doty@gmail.com",
    license='GPL',
    version='0.1.23',
    url='https://github.com/vector/vecash',
    download_url='https://github.com/vector/vecash/tarball/v23',
    description=('Python wrappers for vecash, the vector proof of work'
                 'hashing function'),
    ext_modules=[pyvecash],
)
