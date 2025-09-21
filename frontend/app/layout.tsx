import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { AuthProvider } from '../context/AuthContext';
import Header from '../components/Header';
import './styles.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
    title: 'Minify',
    description: 'Fast and reliable URL minifying service',
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en">
            <body className={inter.className}>
                <AuthProvider>
                    <div >
                        <Header />
                        <main >
                            {children}
                        </main>
                    </div>
                </AuthProvider>
            </body>
        </html>
    );
}
