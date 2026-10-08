'use client';

import { useEffect, useState } from 'react';
import { Search, ArrowUpLeft } from 'lucide-react';
import { Book, books } from '@/lib/api';
import { BookCard } from '@/components/book-card';
import { Hero } from '@/components/landing/hero';
import { StorySection } from '@/components/landing/story-section';
import { ReadingPaths } from '@/components/landing/reading-paths';
import { MembershipCta } from '@/components/landing/membership-cta';

export default function Home() {
  const [items, setItems] = useState<Book[]>([]);
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let active = true;
    setLoading(true);
    const timer = setTimeout(() => {
      books.list(query).then((result) => {
        if (active) { setItems(result); setError(''); }
      }).catch((e) => { if (active) setError(e.message); })
        .finally(() => { if (active) setLoading(false); });
    }, 250);
    return () => { active = false; clearTimeout(timer); };
  }, [query]);

  return <main>
    <Hero />
    <section id="books" className="scroll-mt-24 bg-paper py-20 md:py-28">
      <div className="mx-auto max-w-7xl px-5 md:px-8">
        <div className="scroll-reveal flex flex-col justify-between gap-8 md:flex-row md:items-end">
          <div><p className="section-kicker text-muted">مختارات قارئ / 01</p><h2 className="mt-5 text-4xl font-extrabold leading-tight md:text-6xl">كتب تستحق<br/><span className="font-normal text-muted">أن تتوقف عندها.</span></h2></div>
          <div className="max-w-sm"><p className="leading-relaxed text-muted">من روايات تترك سؤالًا، إلى أفكار تفتح بابًا جديدًا. ابحث عما يشبه فضولك.</p><span className="mt-5 inline-flex items-center gap-2 text-sm font-bold text-ink">اكتشف مجموعتنا <ArrowUpLeft size={17}/></span></div>
        </div>
        <label id="book-search" className="glass-light scroll-reveal mt-12 flex w-full scroll-mt-28 items-center gap-4 px-5 py-4 md:max-w-lg"><Search size={20} strokeWidth={1.5} className="text-muted"/><input aria-label="ابحث عن كتاب أو مؤلف" value={query} onChange={(e) => setQuery(e.target.value)} placeholder="ابحث عن عنوان أو مؤلف..." className="w-full bg-transparent text-base text-ink outline-none placeholder:text-muted/70"/><span className="hidden border-r border-black/15 pr-4 text-xs font-bold tracking-widest text-muted sm:block">SEARCH</span></label>
        {error ? <p role="alert" className="mt-10 rounded-xl border border-black/15 bg-white p-6 text-ink">{error}</p> : loading ? <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3" aria-label="جارٍ تحميل الكتب">{[1,2,3].map((n)=><div key={n} className="aspect-[4/5] animate-pulse rounded-2xl bg-mist"/>)}</div> : items.length ? <div className="mt-10 grid gap-9 sm:grid-cols-2 lg:grid-cols-3">{items.map((book,index)=><BookCard key={book.id} book={book} index={index}/>)}</div> : <p className="mt-10 rounded-xl border border-black/15 bg-white p-10 text-center text-muted">لم نجد كتبًا تطابق بحثك. جرّب عنوانًا أو مؤلفًا آخر.</p>}
      </div>
    </section>
    <StorySection />
    <ReadingPaths />
    <MembershipCta />
  </main>;
}
