"use client";
import { use, useEffect, useState } from "react";
import Link from "next/link";
import { ArrowRight, Star } from "lucide-react";
import { books, library, Book, Review, token } from "@/lib/api";
import { Button } from "@/components/ui/button";
export default function BookPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const [book, setBook] = useState<Book | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [rating, setRating] = useState(5);
  const [body, setBody] = useState("");
  const [userID, setUserID] = useState<number | null>(null);
  const [status, setStatus] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);
  const refresh = () =>
    books.get(id).then((v) => {
      setBook(v.book);
      setReviews(v.reviews);
    });
  useEffect(() => {
    refresh()
      .catch((e) => setMessage(e.message))
      .finally(() => setLoading(false));
    if (token()) {
      import("@/lib/api").then(({ auth, library }) => {
        auth
          .me()
          .then((u) => {
            setUserID(u.id);
            return library.list();
          })
          .then((items) =>
            setStatus(items.find((b) => b.id === Number(id))?.status || ""),
          )
          .catch(() => {});
      });
    }
  }, [id]);
  useEffect(() => {
    const mine = reviews.find((r) => r.user_id === userID);
    if (mine) {
      setRating(mine.rating);
      setBody(mine.body);
    }
  }, [reviews, userID]);
  const saveReview = async () => {
    try {
      await books.review(id, rating, body);
      await refresh();
      setMessage("تم حفظ مراجعتك");
    } catch (e) {
      setMessage((e as Error).message);
    }
  };
  const deleteReview = async () => {
    try {
      await books.deleteReview(id);
      setBody("");
      setRating(5);
      await refresh();
      setMessage("تم حذف المراجعة");
    } catch (e) {
      setMessage((e as Error).message);
    }
  };
  const saveStatus = async (value: string) => {
    try {
      if (value) {
        await library.save(Number(id), value);
      } else {
        await library.remove(Number(id));
      }
      setStatus(value);
      setMessage("تم تحديث مكتبتك");
    } catch (e) {
      setMessage((e as Error).message);
    }
  };
  if (loading)
    return (
      <main className="mx-auto max-w-6xl px-5 py-20">جارٍ تحميل الكتاب...</main>
    );
  if (!book)
    return (
      <main className="mx-auto max-w-6xl px-5 py-20">
        {message || "الكتاب غير موجود"}
      </main>
    );
  const mine = reviews.find((r) => r.user_id === userID);
  return (
    <main className="mx-auto min-h-[70vh] max-w-7xl px-5 py-12 md:px-8 md:py-16">
      <Link
        href="/"
        className="inline-flex items-center gap-2 border-b border-black/30 pb-1 text-sm font-bold text-ink"
      >
        <ArrowRight size={16} />
        العودة للكتب
      </Link>
      <section className="mt-10 grid gap-10 border border-black/10 bg-white p-6 card-shadow md:grid-cols-[330px_1fr] md:gap-16 md:p-12">
        <div className="flex min-h-[420px] items-center justify-center rounded-xl bg-[#e8e8e6] p-8">
          <img
            src={book.cover_url}
            alt={`غلاف ${book.title}`}
            className="editorial-cover max-h-96 max-w-full object-contain drop-shadow-[12px_18px_22px_rgba(0,0,0,.22)]"
          />
        </div>
        <div>
          <span className="section-kicker text-muted">{book.genre}</span>
          <h1 className="mt-5 text-4xl font-extrabold leading-tight md:text-6xl">{book.title}</h1>
          <p className="mt-4 text-lg text-muted">تأليف {book.author}</p>
          <div className="mt-6 flex items-center gap-2">
            <Star className="fill-ink text-ink" />
            <b className="text-xl">
              {book.average_rating ? book.average_rating.toFixed(1) : "—"}
            </b>
            <span className="text-muted">من {book.review_count} مراجعة</span>
          </div>
          <h2 className="mt-10 border-t border-black/10 pt-8 text-lg font-extrabold">عن الكتاب</h2>
          <p className="mt-3 max-w-2xl leading-relaxed text-muted">
            {book.description}
          </p>
          <div className="mt-8 border-t border-black/10 pt-6">
            {userID ? (
              <label className="flex flex-wrap items-center gap-3 text-sm font-bold">
                أضف إلى مكتبتي{" "}
                <select
                  aria-label="تصنيف الكتاب في مكتبتي"
                  className="rounded-xl border border-black/20 bg-paper px-4 py-3"
                  value={status}
                  onChange={(e) => saveStatus(e.target.value)}
                >
                  <option value="">بدون تصنيف</option>
                  <option value="want">أريد قراءته</option>
                  <option value="reading">أقرأه حاليًا</option>
                  <option value="read">قرأته</option>
                </select>
              </label>
            ) : (
              <Link href="/login" className="font-bold text-ink underline underline-offset-4">
                سجّل الدخول لإضافة الكتاب إلى مكتبتك
              </Link>
            )}
          </div>
        </div>
      </section>
      <section className="mt-20 grid gap-10 lg:grid-cols-[1fr_380px]">
        <div>
          <h2 className="text-3xl font-extrabold">
            آراء القراء <span className="text-muted">({reviews.length})</span>
          </h2>
          <div className="mt-6 space-y-4">
            {reviews.length ? (
              reviews.map((r) => (
                <article
                  key={r.id}
                  className="rounded-xl border border-black/10 bg-white p-7"
                >
                  <div className="flex items-center justify-between">
                    <strong>{r.user_name}</strong>
                    <span className="flex items-center gap-1 text-sm font-bold">
                      <Star size={16} className="fill-ink text-ink" />
                      {r.rating}/5
                    </span>
                  </div>
                  <p className="mt-4 leading-relaxed text-muted">{r.body}</p>
                </article>
              ))
            ) : (
              <p className="rounded-xl border border-black/10 bg-white p-7 text-muted">
                كن أول من يراجع هذا الكتاب.
              </p>
            )}
          </div>
        </div>
        <aside className="glass-light h-fit p-7 md:p-8">
          <h2 className="text-xl font-black">
            {mine ? "عدّل مراجعتك" : "شارك رأيك"}
          </h2>
          {userID ? (
            <>
              <p className="mt-4 text-sm font-bold">تقييمك</p>
              <div className="mt-2 flex gap-1" dir="ltr">
                {[1, 2, 3, 4, 5].map((n) => (
                  <button
                    key={n}
                    aria-label={`${n} نجوم`}
                    onClick={() => setRating(n)}
                  >
                    <Star
                      size={27}
                      className={
                        n <= rating ? "fill-ink text-ink" : "text-ink/20"
                      }
                    />
                  </button>
                ))}
              </div>
              <textarea
                aria-label="نص المراجعة"
                value={body}
                onChange={(e) => setBody(e.target.value)}
                placeholder="ما الذي أعجبك في الكتاب؟"
                className="mt-5 h-32 w-full rounded-xl border border-black/20 bg-white/80 p-4 outline-none focus:border-ink"
              />
              <div className="mt-3 flex gap-2">
                <Button onClick={saveReview}>حفظ المراجعة</Button>
                {mine && (
                  <Button variant="outline" onClick={deleteReview}>
                    حذف
                  </Button>
                )}
              </div>
            </>
          ) : (
            <p className="mt-4 text-sm leading-7">
              <Link href="/login" className="font-bold text-ink underline underline-offset-4">
                سجّل الدخول
              </Link>{" "}
              لتكتب مراجعتك.
            </p>
          )}
          {message && (
            <p role="status" className="mt-4 text-sm text-muted">
              {message}
            </p>
          )}
        </aside>
      </section>
    </main>
  );
}
