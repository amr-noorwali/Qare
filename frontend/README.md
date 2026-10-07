> **إلزامي عند العمل على ميزة في الواجهة:** اقرأ ملفها في [`../docs/features/`](../docs/features/README.md) قبل التنفيذ، وحدّثه في نفس المهمة بعد أي تغيير. راجع [`../AGENTS.md`](../AGENTS.md).

# واجهة قَرأ

واجهة عربية مبنية بـ Next.js وReact وTypeScript وTailwind CSS، مع مكونات UI بأسلوب shadcn/ui داخل `components/ui`.

## التشغيل

```powershell
Copy-Item .env.example .env.local
npm install
npm run dev
```

افتح http://localhost:3000 بعد تشغيل الباك إند على المنفذ 8080. `NEXT_PUBLIC_API_URL` هو عنوان REST API الكامل مع `/api`. للتحقق: `npm run lint` ثم `npm run build`.

تُحفظ جلسة الدخول في `localStorage` في هذه النسخة الأولية. للاستخدام الإنتاجي، انقل الجلسات إلى ملفات تعريف ارتباط `HttpOnly` مع HTTPS.
