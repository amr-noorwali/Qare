import Link from 'next/link';
import { ArrowUpLeft, BookOpen, MessageSquareQuote, LibraryBig } from 'lucide-react';

const values = [
  { icon: BookOpen, title: 'اكتشاف', description: 'ابدأ من كتاب يثير فضولك.' },
  { icon: MessageSquareQuote, title: 'انطباع', description: 'أضف صوتك إلى الحديث عنه.' },
  { icon: LibraryBig, title: 'أثر', description: 'احتفظ بمسار قراءتك في مكتبتك.' },
];

export function StorySection() {
  return (
    <section id="story" className="bg-[#ececea] py-20 md:py-28">
      <div className="mx-auto grid max-w-7xl gap-10 px-5 md:grid-cols-[.95fr_1.05fr] md:items-center md:gap-20 md:px-8">
        <div className="relative min-h-[470px] overflow-hidden bg-ink md:min-h-[680px]">
          <img src="/editorial-reading.webp" alt="صفحات كتاب تُقلب في مساحة قراءة هادئة" loading="lazy" className="absolute inset-0 h-full w-full object-cover grayscale" />
          <div className="glass-dark absolute bottom-5 left-5 right-5 px-6 py-5 text-white md:bottom-8 md:left-8 md:right-auto md:max-w-[290px]">
            <p className="section-kicker text-white/55">من الصفحة إلى الذاكرة</p>
            <p className="mt-2 text-lg font-bold leading-7">كل قراءة تترك أثرًا يستحق أن يُحفظ.</p>
          </div>
        </div>
        <div className="py-2 md:py-10">
          <p className="section-kicker text-muted">حكاية قارئ</p>
          <h2 className="mt-7 max-w-lg text-4xl font-extrabold leading-[1.25] tracking-tight md:text-6xl">أكثر من كتاب.<br/><span className="font-normal text-muted">مساحة لما يتركه فيك.</span></h2>
          <p className="mt-7 max-w-xl text-base leading-9 text-muted">نؤمن بأن الكتاب لا ينتهي عند صفحته الأخيرة. قارئ مساحة بسيطة لتكتشف ما تقرأ، تشارك رأيك بصدق، وتعود إلى الكتب التي صنعت جزءًا من رحلتك.</p>
          <div className="mt-9 divide-y divide-black/15 border-y border-black/15">
            {values.map(({icon: Icon,title,description}, index) => <div key={title} className="flex items-center gap-5 py-5"><span className="flex h-12 w-12 shrink-0 items-center justify-center border border-black/20"><Icon size={20} strokeWidth={1.5}/></span><div className="flex-1"><h3 className="font-extrabold">{title}</h3><p className="mt-1 text-sm text-muted">{description}</p></div><span className="text-xs text-muted">0{index+1}</span></div>)}
          </div>
          <Link href="/library" className="mt-9 inline-flex items-center gap-3 border-b border-ink pb-2 text-sm font-bold">اكتشف مكتبتك <ArrowUpLeft size={17}/></Link>
        </div>
      </div>
    </section>
  );
}
