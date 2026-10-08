import Link from 'next/link';
import { ArrowUpLeft } from 'lucide-react';

export function MembershipCta() {
  return <section className="bg-paper px-5 pb-8 md:px-8 md:pb-14"><div className="scroll-reveal relative mx-auto max-w-7xl overflow-hidden rounded-2xl bg-ink px-6 py-20 text-center text-white md:px-16 md:py-28"><div className="absolute -right-24 -top-40 h-96 w-96 rounded-full bg-white/10 blur-[90px]"/><div className="absolute -bottom-60 left-0 h-96 w-96 rounded-full bg-white/10 blur-[100px]"/><div className="glass-dark relative mx-auto max-w-4xl px-6 py-14 md:px-16 md:py-20"><p className="section-kicker text-white/75">مساحة القارئ</p><h2 className="mt-6 text-4xl font-extrabold leading-tight md:text-6xl">انضم إلى عالم<br/>تبدأ فيه الحكاية منك.</h2><p className="mx-auto mt-6 max-w-xl text-base leading-relaxed text-white/80">أنشئ حسابك لتكتب مراجعاتك، ترتب كتبك، وتبني مكتبة تشبه ذائقتك.</p><Link href="/register" className="interactive-lift mt-9 inline-flex items-center gap-6 rounded-xl bg-white px-8 py-4 text-sm font-extrabold text-ink transition hover:bg-white/85">أنشئ حسابك <ArrowUpLeft size={18}/></Link></div></div></section>;
}
