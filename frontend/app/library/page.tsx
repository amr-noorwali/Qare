"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { BookCard } from "@/components/book-card";
import { library, LibraryBook, token } from "@/lib/api";
const tabs = [
  ["want", "أريد قراءته"],
  ["reading", "أقرأه حاليًا"],
  ["read", "قرأته"],
] as const;
export default function LibraryPage() {
  const [items, setItems] = useState<LibraryBook[]>([]);
  const [tab, setTab] = useState<"want" | "reading" | "read">("want");
  const [error, setError] = useState("");
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  useEffect(() => {
    if (!token()) {
      setIsAuthenticated(false);
      return;
    }
    setIsAuthenticated(true);
    library
      .list()
      .then(setItems)
      .catch((e) => setError(e.message));
  }, []);
  return (
    <main className="mx-auto min-h-[70vh] max-w-7xl px-5 py-16 md:px-8 md:py-24">
      <span className="section-kicker text-muted">مساحتي الخاصة / 02</span>
      <h1 className="page-enter mt-5 text-5xl font-extrabold md:text-6xl">مكتبتي<span className="text-muted">.</span></h1>
      <p className="mt-5 text-lg text-muted">كل الكتب التي اخترتها، في مكان واحد.</p>
      {isAuthenticated === null ? (
        <div className="glass-light mt-12 p-12 text-center text-muted">
          جارٍ تحميل مكتبتك...
        </div>
      ) : !isAuthenticated ? (
        <div className="glass-light mt-12 p-12 text-center">
          <p>سجّل الدخول لتبدأ بتنظيم مكتبتك.</p>
          <Link
            href="/login"
            className="interactive-lift mt-6 inline-block rounded-xl border border-ink bg-ink px-7 py-3 font-bold text-white transition hover:bg-white hover:text-ink"
          >
            تسجيل الدخول
          </Link>
        </div>
      ) : (
        <>
          <div className="mt-14 flex gap-3 overflow-x-auto border-b border-black/15 pb-4">
            {tabs.map(([key, label]) => (
              <button
                key={key}
                onClick={() => setTab(key)}
                className={`interactive-lift whitespace-nowrap rounded-xl border px-5 py-3 text-sm font-bold transition ${tab === key ? "border-ink bg-ink text-white" : "border-black/15 bg-white/70 text-ink hover:border-ink"}`}
              >
                {label} ({items.filter((b) => b.status === key).length})
              </button>
            ))}
          </div>
          {error ? (
            <p className="mt-8 text-red-700">{error}</p>
          ) : items.filter((b) => b.status === tab).length ? (
            <div className="mt-10 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
              {items
                .filter((b) => b.status === tab)
                .map((b) => (
                  <BookCard key={b.id} book={b} />
                ))}
            </div>
          ) : (
            <div className="glass-light mt-10 p-12 text-center text-muted">
              لا توجد كتب في هذا القسم بعد.{" "}
              <Link href="/" className="font-bold underline">
                استكشف الكتب
              </Link>
            </div>
          )}
        </>
      )}
    </main>
  );
}
