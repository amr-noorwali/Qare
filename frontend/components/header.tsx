'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { Menu, X } from 'lucide-react';
import { auth, User } from '@/lib/api';

export function Header() {
  const [user, setUser] = useState<User | null>(null);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (!localStorage.getItem('qare_token')) return;
    auth.me().then(setUser).catch(() => {
      localStorage.removeItem('qare_token');
      setUser(null);
    });
  }, []);

  const logout = () => {
    localStorage.removeItem('qare_token');
    setUser(null);
    location.href = '/';
  };

  return (
    <header className="sticky top-0 z-50 border-b border-black/10 bg-white/75 backdrop-blur-2xl">
      <div className="mx-auto flex h-[76px] max-w-7xl items-center justify-between px-5 md:px-8">
        <Link href="/" className="flex items-center gap-3 text-2xl font-extrabold tracking-tight text-ink" aria-label="قارئ، الصفحة الرئيسية">
          <img src="/qare-logo.svg" alt="" width={36} height={36} className="h-9 w-9 object-contain" />
          <span>قارئ</span>
        </Link>
        <nav className="hidden items-center gap-10 text-[13px] font-bold text-ink/75 md:flex" aria-label="التنقل الرئيسي">
          <Link href="/#books" className="transition hover:text-black">استكشف الكتب</Link>
          <Link href="/#story" className="transition hover:text-black">عن قارئ</Link>
          <Link href="/library" className="transition hover:text-black">مكتبتي</Link>
        </nav>
        <div className="hidden items-center gap-5 md:flex">
          {user ? <><span className="text-sm text-muted">مرحبًا، {user.name}</span><button onClick={logout} className="border-b border-ink pb-1 text-sm font-bold">تسجيل الخروج</button></> : <><Link href="/login" className="text-sm font-bold">تسجيل الدخول</Link><Link href="/register" className="border border-ink bg-ink px-5 py-2.5 text-sm font-bold text-white transition hover:bg-white hover:text-ink">انضم إلى قارئ</Link></>}
        </div>
        <button className="inline-flex h-10 w-10 items-center justify-center border border-black/15 md:hidden" aria-label={open ? 'إغلاق القائمة' : 'فتح القائمة'} aria-expanded={open} onClick={() => setOpen(!open)}>{open ? <X size={21}/> : <Menu size={21}/>}</button>
      </div>
      {open && <nav className="glass-light absolute inset-x-3 top-[82px] z-50 flex flex-col gap-1 p-4 text-sm font-bold md:hidden" aria-label="تنقل الجوال">
        <Link onClick={() => setOpen(false)} href="/#books" className="px-3 py-3">استكشف الكتب</Link>
        <Link onClick={() => setOpen(false)} href="/#story" className="px-3 py-3">عن قارئ</Link>
        <Link onClick={() => setOpen(false)} href="/library" className="px-3 py-3">مكتبتي</Link>
        <div className="my-2 border-t border-black/10" />
        {user ? <button className="px-3 py-3 text-right" onClick={logout}>تسجيل الخروج</button> : <><Link onClick={() => setOpen(false)} href="/login" className="px-3 py-3">تسجيل الدخول</Link><Link onClick={() => setOpen(false)} href="/register" className="bg-ink px-3 py-3 text-white">انضم إلى قارئ</Link></>}
      </nav>}
    </header>
  );
}
