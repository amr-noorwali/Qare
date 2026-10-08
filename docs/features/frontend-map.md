> **إلزامي للفرونت إند:** اقرأ هذا الملف قبل أي مهمة تمس بنية الواجهة أو توزيع ملفات ميزة فيها. بعد التغيير، حدّث خريطة الملفات والحالة وسجل التغييرات هنا، وحدّث ملف الميزة المعنية في نفس المهمة.

# خريطة الفرونت إند: أين يوجد تنفيذ كل ميزة؟

توجد ملفات توثيق الفرونت نفسها في [frontend/docs/features/](../../frontend/docs/features/README.md)، بينما هذه الخريطة تشرح توزيع كود Next.js.

## لماذا تبدو الميزات كأنها صفحات فقط؟

المجلد `frontend/app/` يستخدم **Next.js App Router**؛ لذلك يحتوي تعريف المسارات (`page.tsx`) وملف التخطيط العام. لم يُبن المشروع بمجلد `features/` مستقل لكل ميزة. ملف الصفحة قد يحتوي حاليًا على بعض حالة الواجهة وتفاعلاتها، بينما المكونات القابلة لإعادة الاستخدام في `frontend/components/` والاتصال بالباك إند في `frontend/lib/api.ts`. التوثيق تحت `docs/features/` يصنف **السلوك حسب الميزة**، وليس هيكل المجلدات داخل الفرونت.

## خريطة الطبقات

| المكان | مسؤوليته الحالية | أمثلة |
|---|---|---|
| `app/` | المسارات، تركيب الشاشة، حالة الصفحة وتفاعلاتها | `app/page.tsx`, `app/library/page.tsx`, `app/books/[id]/page.tsx` |
| `components/` | واجهة مشتركة أو قابلة لإعادة الاستخدام | `header.tsx`, `book-card.tsx`, `auth-form.tsx` |
| `components/landing/` | أقسام الصفحة الرئيسية التحريرية، منفصلة عن البحث وعرض الكتب | `hero.tsx`, `story-section.tsx`, `reading-paths.tsx`, `membership-cta.tsx` |
| `components/ui/` | مكونات UI عامة | `button.tsx` |
| `lib/api.ts` | أنواع بيانات الواجهة وعميل REST ودوال `auth`, `books`, `library` | `api()`, `books.list()`, `library.save()` |
| `lib/utils.ts` | دمج أسماء CSS | `cn()` |
| `app/globals.css`, `tailwind.config.ts` | خط Tajawal ومقياس أحجام النص وتباعد الأسطر والألوان وتأثيرات الزجاج | أحجام `text-xs` إلى `text-7xl`، وألوان `ink/paper/mist/line/muted` |
| `public/`, `app/icon.svg` | أصول الشعار والأيقونة وصور الهيرو والقصة وخلفية الحساب | `qare-logo.svg`, `hero-book-stack.png`, `story-hands-book.png`, `editorial-library.webp` |

## خريطة الميزات في الفرونت

| الميزة | الصفحة | المكون المشترك | خدمة REST في `lib/api.ts` | حالة/تفاعل الصفحة |
|---|---|---|---|---|
| [التسجيل والدخول](authentication.md) | `app/register/page.tsx`, `app/login/page.tsx` | `components/auth-form.tsx`, `components/header.tsx` | `auth.register`, `auth.login`, `auth.me`, `token` | حقول النموذج، الأخطاء، التحميل، حفظ الرمز والتنقل |
| [الكتب والبحث](catalog-search.md) | `app/page.tsx` | `components/book-card.tsx` | `books.list` | نص البحث، تأخير 250ms، التحميل والنتائج |
| [تفاصيل الكتاب](book-details.md) | `app/books/[id]/page.tsx` | `components/ui/button.tsx` | `books.get`, `auth.me`, `library.list` | تحميل التفاصيل، تحديد المستخدم ومراجعته وتصنيف الكتاب |
| [المراجعات](reviews-ratings.md) | جزء من `app/books/[id]/page.tsx` | `components/ui/button.tsx` | `books.review`, `books.deleteReview` | اختيار النجوم، نص المراجعة، الحفظ والحذف |
| [المكتبة](personal-library.md) | `app/library/page.tsx`، واختيار التصنيف في `app/books/[id]/page.tsx` | `components/book-card.tsx` | `library.list`, `library.save`, `library.remove` | التبويبات، العدادات، حالة الجلسة، تغيير التصنيف |
| [الهوية البصرية](visual-identity.md) | `app/layout.tsx`, `app/globals.css` | `components/header.tsx`, `components/ui/button.tsx` | لا يوجد اتصال API للهوية | RTL، Tajawal، الشعار، الألوان |

## تدفق مثال: البحث

`app/page.tsx` يحتفظ بقيمة البحث → يستدعي `books.list(q)` من `lib/api.ts` → Go يعيد الكتب → `book-card.tsx` يعرض كل كتاب → النقر ينتقل إلى `app/books/[id]/page.tsx`.

## الحالة والحدود

- هذا فصل بسيط مناسب للـMVP، لكنه **ليس تنظيمًا لكل ميزة داخل مجلد مستقل**. إذا كبر المشروع، يمكن نقل حالة الصفحة وعملياتها إلى hooks أو وحدات ميزات دون تغيير مسارات `app/`.
- `lib/api.ts` هو طبقة خدمات الاتصال الحالية؛ لا يوجد مجلد `services/` منفصل في الفرونت.
- صفحة تفاصيل الكتاب تحمل أكثر من مسؤولية: التفاصيل والمراجعات والتحكم في المكتبة. تقسيمها إلى مكونات أصغر تحسين صيانة مستقبلي، وليس مطبقًا حاليًا.

## سجل التغييرات

- 2026-10-08: استُبدل أصل قسم القصة بـ`story-hands-book.png` وحُذف `editorial-reading.webp` غير المستخدم؛ بقي توزيع المكونات والـAPI كما هو.
- 2026-10-08: أضيف مقياس خط مركزي في `tailwind.config.ts` وحجم النص الأساسي في `app/globals.css`، وطُبّق على الصفحات والمكونات دون تغيير مواقعها أو REST API.
- 2026-10-08: استُبدل أصل الهيرو السابق بصورة كومة الكتب المقدمة من المستخدم داخل `public/`؛ لم يتغير توزيع مكونات الصفحة.
- 2026-10-08: أضيفت رسمة الهيرو ثنائية اللون إلى خريطة الأصول؛ بقيت أقسام الهبوط منفصلة في `components/landing/` والبحث في `app/page.tsx`.
- 2026-10-08: أضيفت مكونات الهبوط المستقلة وصور WebP التحريرية، ووُثّقت لوحة الأبيض والأسود وتأثيرات الزجاج بعد إعادة تصميم الواجهة؛ بقي البحث والمصادقة والمراجعات والمكتبة على عميل REST نفسه.
- 2026-10-08: رُبطت الخريطة بفهرس توثيق الميزات داخل مشروع الفرونت.
- 2026-10-08: أضيفت خريطة صريحة للصفحات والمكونات وخدمة API وحالة كل ميزة لتوضيح أن التنفيذ لا يقتصر على `app/`.
