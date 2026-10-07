import './globals.css';
import { Header } from '@/components/header';
export const metadata={title:'قَرأ | مساحة للكتب والقراء',description:'اكتشف الكتب، شارك رأيك، ونظّم مكتبتك الشخصية'};
export default function RootLayout({children}:{children:React.ReactNode}){return <html lang="ar" dir="rtl"><body className="font-sans antialiased"><Header/>{children}<footer className="mt-20 border-t border-moss/10 py-8 text-center text-sm text-moss/60">قارئ © ٢٠٢٦ · لكل كتاب حكاية</footer></body></html>}
