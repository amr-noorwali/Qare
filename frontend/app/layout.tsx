import Link from 'next/link';
import './globals.css';
import { Header } from '@/components/header';

export const metadata = {
  title: 'قارئ | مساحة للكتب التي تبقى',
  description: 'اكتشف كتبًا تستحق وقتك، شارك رأيك، وابنِ مكتبتك الشخصية في قارئ.',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return <html lang="ar" dir="rtl"><body className="font-sans antialiased">
    <Header />
    {children}
    <footer className="bg-ink text-white">
      <div className="mx-auto grid max-w-7xl gap-12 px-5 py-16 md:grid-cols-[1.5fr_1fr_1fr] md:px-8 md:py-20">
        <div><Link href="/" className="inline-flex items-center gap-3 text-2xl font-extrabold"><img src="/qare-logo.svg" alt="" width={38} height={38} className="h-9 w-9 invert"/>قارئ</Link><p className="mt-5 max-w-sm text-sm leading-relaxed text-white/70">مساحة للكتب والقراء. اكتشف ما يستحق أن يُقرأ، واحتفظ بأثر كل كتاب.</p></div>
        <div><h2 className="section-kicker text-white/70">استكشف</h2><div className="mt-6 flex flex-col gap-4 text-sm"><Link href="/#books" className="hover:text-white/70">الكتب</Link><Link href="/library" className="hover:text-white/70">مكتبتي</Link><Link href="/#story" className="hover:text-white/70">حكاية قارئ</Link></div></div>
        <div><h2 className="section-kicker text-white/70">حسابك</h2><div className="mt-6 flex flex-col gap-4 text-sm"><Link href="/register" className="hover:text-white/70">إنشاء حساب</Link><Link href="/login" className="hover:text-white/70">تسجيل الدخول</Link></div></div>
      </div>
      <div className="mx-auto flex max-w-7xl flex-col justify-between gap-3 border-t border-white/15 px-5 py-6 text-xs text-white/70 md:flex-row md:px-8"><span>© ٢٠٢٦ قارئ. لكل كتاب حكاية.</span><span>صُنعت بحب للقراءة.</span></div>
    </footer>
  </body></html>;
}
