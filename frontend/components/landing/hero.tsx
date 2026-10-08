import Image from 'next/image';
import Link from 'next/link';
import { ArrowUpLeft, Search } from 'lucide-react';

export function Hero() {
  return (
    <section className="bg-white text-black" aria-labelledby="hero-title">
      <div className="mx-auto max-w-7xl px-5 pt-10 md:px-8 md:pt-14">
        <div className="grid items-center gap-10 lg:grid-cols-[1fr_1.05fr] lg:gap-16">
          <div className="order-1 py-4 lg:py-12">
            <div className="flex items-center gap-4 text-sm font-bold md:text-base"><span className="h-px w-10 bg-black" /><span>لكل قارئ، حكاية.</span></div>
            <h1 id="hero-title" className="mt-7 max-w-[630px] text-[clamp(3.25rem,6.6vw,7rem)] font-extrabold leading-[1.08] tracking-tight">بين دفّتين،<br />عالم ينتظرك.</h1>
            <p className="mt-7 max-w-xl text-lg leading-9 md:text-2xl md:leading-[1.8]">اكتشف كتبًا تأخذك أبعد، وحكايات تبقى معك.</p>
            <div className="mt-9 flex flex-wrap gap-3">
              <a href="#books" className="inline-flex min-h-14 items-center justify-center gap-5 border border-black bg-black px-7 py-3 text-sm font-bold text-white transition-colors hover:bg-white hover:text-black">اكتشف الكتب <ArrowUpLeft size={18} aria-hidden="true" /></a>
              <a href="#book-search" className="inline-flex min-h-14 items-center justify-center gap-4 border border-black bg-white px-7 py-3 text-sm font-bold text-black transition-colors hover:bg-black hover:text-white">ابحث عن كتاب <Search size={18} aria-hidden="true" /></a>
            </div>
            <p className="mt-7 text-sm leading-7">احتفظ بمراجعاتك وكتبك المفضلة في <Link href="/register" className="font-extrabold underline underline-offset-4">مكتبتك الخاصة</Link>.</p>
          </div>
          <div className="order-2 h-[420px] overflow-hidden bg-black sm:h-[550px] lg:h-[650px]">
            <Image src="/hero-book-stack.png" alt="رسم أبيض وأسود لكومة كتب متراصة" width={1214} height={1295} priority unoptimized className="h-full w-full object-contain" />
          </div>
        </div>
        <div className="mt-10 flex items-center gap-5 border-y border-black py-5 md:mt-14">
          <div className="h-px flex-1 bg-black" />
          <p className="flex flex-wrap justify-center gap-x-4 gap-y-1 text-xs font-bold sm:text-sm" aria-label="أنواع الكتب"><span>روايات</span><span aria-hidden="true">•</span><span>أدب</span><span aria-hidden="true">•</span><span>فكر</span><span aria-hidden="true">•</span><span>تاريخ</span><span aria-hidden="true">•</span><span>تطوير الذات</span></p>
          <div className="h-px flex-1 bg-black" />
        </div>
      </div>
    </section>
  );
}
