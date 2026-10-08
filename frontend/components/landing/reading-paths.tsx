import Link from 'next/link';
import { ArrowUpLeft } from 'lucide-react';

const paths = [
  { number: '01', title: 'أريد قراءته', note: 'احتفظ بالكتب التي تنتظر لحظتها.' },
  { number: '02', title: 'أقرأه حاليًا', note: 'تابع ما يصاحب أيامك الآن.' },
  { number: '03', title: 'قرأته', note: 'عد إلى الحكايات التي تركت أثرًا.' },
];

export function ReadingPaths() {
  return <section className="bg-paper py-20 md:py-28"><div className="mx-auto max-w-7xl px-5 md:px-8"><div className="flex flex-col justify-between gap-5 md:flex-row md:items-end"><div><p className="section-kicker text-muted">مكتبتك الخاصة</p><h2 className="mt-5 text-4xl font-extrabold md:text-6xl">لكل كتاب مكان.</h2></div><p className="max-w-md leading-relaxed text-muted">رتّب رحلتك القرائية بثلاث حالات واضحة، وانتقل بين الكتب كما تنتقل بين فصول حكايتك.</p></div><div className="mt-12 grid gap-4 md:grid-cols-3">{paths.map((path,index)=><Link key={path.number} href="/library" className={`group relative flex min-h-[285px] flex-col justify-between overflow-hidden rounded-2xl border p-7 transition duration-300 hover:-translate-y-1 ${index===1?'border-ink bg-ink text-white':'border-black/15 bg-white text-ink'}`}><span className={`text-xs font-bold ${index===1?'text-white/70':'text-muted'}`}>{path.number} / مسار القراءة</span><div><h3 className="text-3xl font-extrabold">{path.title}</h3><p className={`mt-3 text-sm leading-relaxed ${index===1?'text-white/75':'text-muted'}`}>{path.note}</p></div><ArrowUpLeft size={23} className="absolute left-7 top-7 transition group-hover:-translate-x-1 group-hover:-translate-y-1"/></Link>)}</div></div></section>;
}
