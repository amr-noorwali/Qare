import Link from 'next/link';
import { ArrowUpLeft, Star } from 'lucide-react';
import { Book } from '@/lib/api';

export function BookCard({ book, index }: { book: Book; index?: number }) {
  return (
    <Link href={`/books/${book.id}`} className="group block min-w-0" aria-label={`تفاصيل ${book.title}`}>
      <div className="relative flex aspect-[4/5] items-center justify-center overflow-hidden bg-[#e8e8e6]">
        <div className="absolute inset-4 border border-white/55" />
        <img src={book.cover_url} alt={`غلاف ${book.title}`} loading="lazy" className="editorial-cover relative z-10 h-[75%] max-w-[72%] object-contain drop-shadow-[10px_14px_20px_rgba(0,0,0,.23)]" />
        <span className="glass-light absolute bottom-3 left-3 z-20 flex h-10 w-10 items-center justify-center text-ink transition group-hover:bg-ink group-hover:text-white"><ArrowUpLeft size={18}/></span>
        {index !== undefined && <span className="absolute right-5 top-5 z-20 text-xs font-bold tracking-[.2em] text-ink/55">{String(index + 1).padStart(2, '0')}</span>}
      </div>
      <div className="flex items-start justify-between gap-3 border-b border-black/10 py-5">
        <div className="min-w-0">
          <p className="section-kicker text-muted">{book.genre}</p>
          <h3 className="mt-2 truncate text-xl font-extrabold text-ink">{book.title}</h3>
          <p className="mt-1 text-sm text-muted">{book.author}</p>
        </div>
        <div className="flex shrink-0 items-center gap-1 pt-1 text-sm font-bold text-ink"><Star size={14} className="fill-ink" />{book.average_rating ? book.average_rating.toFixed(1) : '—'}</div>
      </div>
    </Link>
  );
}
