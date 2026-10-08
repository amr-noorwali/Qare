> **إلزامي عند العمل على ميزة في الواجهة:** اقرأ ملفها في [`docs/features/`](docs/features/README.md) والملف الشامل المناظر في [`../docs/features/`](../docs/features/README.md) قبل التنفيذ، وحدّث الاثنين في نفس المهمة بعد أي تغيير. راجع [تعليمات الفرونت](AGENTS.md).

# واجهة قارئ

واجهة عربية مبنية بـ Next.js وReact وTypeScript وTailwind CSS، مع مكونات UI بأسلوب shadcn/ui داخل `components/ui`.

توثيق كل ميزة داخل مشروع الواجهة في [frontend/docs/features](docs/features/README.md)، وخريطة توزيع الصفحات والمكونات وخدمات REST في [خريطة الفرونت إند](../docs/features/frontend-map.md). ملفات `app/` تحدد المسارات؛ التنفيذ لا يقتصر عليها.

## التصميم

الواجهة عربية RTL بخط Tajawal وهوية أحادية اللون. رموز الألوان والتدرجات والزجاج في `tailwind.config.ts` و`app/globals.css`، وأقسام الهبوط في `components/landing/`. صور الهبوط الأصلية بصيغة WebP داخل `public/`. روابط العضوية تقود إلى التسجيل الحالي؛ لا توجد خدمة اشتراك مدفوع في هذا الـMVP.

## التشغيل

```powershell
Copy-Item .env.example .env.local
npm install
npm run dev
```

افتح http://localhost:3000 بعد تشغيل الباك إند على المنفذ 8080. `NEXT_PUBLIC_API_URL` هو عنوان REST API الكامل مع `/api`. للتحقق: `npm run lint` ثم `npm run build`.

تُحفظ جلسة الدخول في `localStorage` في هذه النسخة الأولية. للاستخدام الإنتاجي، انقل الجلسات إلى ملفات تعريف ارتباط `HttpOnly` مع HTTPS.
