#![crate_type="staticlib"]
#![feature(start)]
#![feature(lang_items)]
#![no_std]

extern crate libc;

extern {
    fn __gray_main();
}

#[start]
fn start(_argc: isize, _argv: *const *const u8) -> isize {
    unsafe {
        __gray_main();
    }
    0
}

#[lang = "eh_personality"]
#[no_mangle]
pub extern fn eh_personality() {
}

#[lang = "panic_fmt"]
#[no_mangle]
pub extern fn rust_begin_panic(_msg: core::fmt::Arguments,
                               _file: &'static str,
                               _line: u32) -> ! {
    loop {}
}
