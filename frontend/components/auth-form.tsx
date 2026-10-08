'use client';

import { FormEvent, useState } from 'react';
import Link from 'next/link';
import { auth } from '@/lib/api';
import { Button } from './ui/button';

export function AuthForm({ mode }: { mode: 'login' | 'register' }) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setBusy(true);
    setError('');
    try {
      const result = mode === 'register' ? await auth.register(name, email, password) : await auth.login(email, password);
      localStorage.setItem('qare_token', result.token);
      location.href = '/library';
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  };

  return <main className="relative min-h-[calc(100vh-76px)] overflow-hidden bg-[#eaeae8] px-5 py-16 md:py-24">
    <div className="absolute inset-0 bg-[url('/editorial-library.webp')] bg-cover bg-center opacity-[.17] grayscale" />
    <div className="absolute inset-0 bg-gradient-to-b from-white/65 via-white/55 to-white/80" />
    <div className="glass-light relative mx-auto max-w-lg p-7 md:p-12">
      <Link href="/" className="inline-flex items-center gap-3 text-xl font-extrabold"><img src="/qare-logo.svg" alt="" width={40} height={40} className="h-10 w-10 object-contain"/>قارئ</Link>
      <p className="section-kicker mt-12 text-muted">مساحة القارئ</p>
      <h1 className="mt-4 text-4xl font-extrabold leading-tight">{mode === 'register' ? 'ابدأ حكايتك هنا.' : 'مرحبًا بعودتك.'}</h1>
      <p className="mt-3 leading-8 text-muted">{mode === 'register' ? 'حساب واحد يجمع مراجعاتك وكتبك ورحلتك القرائية.' : 'أكمل رحلتك بين الكتب التي اخترتها.'}</p>
      <form onSubmit={submit} className="mt-10 space-y-6">
        {mode === 'register' && <label className="block text-sm font-bold">الاسم<input required minLength={2} value={name} onChange={(e) => setName(e.target.value)} className="mt-2 w-full border border-black/20 bg-white/75 px-4 py-3.5 outline-none transition focus:border-ink" /></label>}
        <label className="block text-sm font-bold">البريد الإلكتروني<input required type="email" value={email} onChange={(e) => setEmail(e.target.value)} className="mt-2 w-full border border-black/20 bg-white/75 px-4 py-3.5 outline-none transition focus:border-ink" /></label>
        <label className="block text-sm font-bold">كلمة المرور<input required minLength={8} type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="mt-2 w-full border border-black/20 bg-white/75 px-4 py-3.5 outline-none transition focus:border-ink" /></label>
        {error && <p role="alert" className="border-r-2 border-ink bg-white/70 px-4 py-3 text-sm text-ink">{error}</p>}
        <Button disabled={busy} className="w-full" type="submit">{busy ? 'جارٍ المعالجة...' : mode === 'register' ? 'إنشاء حساب' : 'تسجيل الدخول'}</Button>
      </form>
      <p className="mt-8 border-t border-black/10 pt-6 text-center text-sm text-muted">{mode === 'register' ? 'لديك حساب؟' : 'جديد هنا؟'} <Link className="font-extrabold text-ink underline underline-offset-4" href={mode === 'register' ? '/login' : '/register'}>{mode === 'register' ? 'سجّل الدخول' : 'أنشئ حسابًا'}</Link></p>
    </div>
  </main>;
}
