-- ========================================================
-- KPM ACADEMY - OFFICIAL STANDARDIZED CATALOG SEED
-- Membersihkan data dummy dan mengisi katalog resmi KPM MIPA
-- ========================================================

-- 1. Bersihkan tabel katalog lama
DELETE FROM banners;
DELETE FROM packages;
DELETE FROM live_classes;
DELETE FROM videos;
DELETE FROM modules;
DELETE FROM vouchers;

-- 2. SPANDUK BANNER RESMI BERANDA MOBILE
INSERT INTO banners (id, tag, sub_tag, title, subtitle, image_url, route, is_active, sort_order, created_at, updated_at) VALUES
('b1111111-1111-1111-1111-111111111111', 'PROGRAM UNGGULAN', 'Metode KPM', 'Matematika Nalaria Realistik (MNR)', 'Metode berpikir kritis dan pemecahan masalah matematika kontekstual khas KPM', 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=900&auto=format&fit=crop&q=80', '/package_list', 1, 1, NOW(), NOW()),
('b2222222-2222-2222-2222-222222222222', 'TRYOUT CBT', 'Skor Otomatis', 'Simulasi Ujian & Tryout CBT Online', 'Ribuan bank soal berstandar olimpiade dengan evaluasi peringkat nasional', 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=900&auto=format&fit=crop&q=80', '/package_list', 1, 2, NOW(), NOW()),
('b3333333-3333-3333-3333-333333333333', 'KELAS LIVE', 'Tatap Muka Daring', 'Sesi Interaktif Live Class Bersama Master Tutor', 'Bimbingan tatap muka daring melalui Zoom & Google Meet berjadwal', 'https://images.unsplash.com/photo-1588072432836-e10032774350?w=900&auto=format&fit=crop&q=80', '/live_classes', 1, 3, NOW(), NOW()),
('b4444444-4444-4444-4444-444444444444', 'E-BOOK & MODUL', 'Unduh PDF', 'Modul & Bank Soal Digital format PDF', 'Arsip materi belajar, rangkuman rumus, dan pembahasan soal terstruktur', 'https://images.unsplash.com/photo-1524178232363-1fb2b075b655?w=900&auto=format&fit=crop&q=80', '/package_list', 1, 4, NOW(), NOW());

-- 3. PAKET BELAJAR & TRYOUT RESMI
INSERT INTO packages (id, title, description, thumbnail, kelas, jenjang, price, discount_price, is_discount_active, membership_duration_days, is_active, created_at, updated_at) VALUES
('p1111111-1111-1111-1111-111111111111', 'Matematika Nalaria Realistik (MNR) SD', 'Pembinaan logika bernalar, pemahaman pecahan, deret angka, dan penyelesaian soal cerita aplikatif untuk siswa sekolah dasar.', 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=600&auto=format&fit=crop&q=80', 'Kelas 4-6 SD', 'SD', 150000.00, 99000.00, 1, 180, 1, NOW(), NOW()),
('p2222222-2222-2222-2222-222222222222', 'Intensif Juara Olimpiade Matematika (OSN / KMNR) SMP', 'Bedah tipe soal olimpiade sains dan matematika tingkat kabupaten hingga nasional lengkap dengan trik pemecahan masalah cepat.', 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=600&auto=format&fit=crop&q=80', 'Kelas 7-9 SMP', 'SMP', 250000.00, 149000.00, 1, 365, 1, NOW(), NOW()),
('p3333333-3333-3333-3333-333333333333', 'Mastering Aljabar & Geometri Analitik SMP', 'Pendalaman materi SPLDV, fungsi kuadrat, phytagoras, lingkaran, serta bangun ruang dengan metode visual interaktif.', 'https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=600&auto=format&fit=crop&q=80', 'Kelas 8-9 SMP', 'SMP', 175000.00, 119000.00, 1, 180, 1, NOW(), NOW()),
('p4444444-4444-4444-4444-444444444444', 'Persiapan UTBK-SNBT Penalaran Matematika & TPS SMA', 'Simulasi tes skolastik, pemahaman bacaan kuantitatif, dan penalaran matematika terstandar seleksi masuk PTN.', 'https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=600&auto=format&fit=crop&q=80', 'Kelas 10-12 SMA', 'SMA', 350000.00, 199000.00, 1, 365, 1, NOW(), NOW()),
('p5555555-5555-5555-5555-555555555555', 'Sains & IPA Terpadu Eksploratif SD-SMP', 'Pembahasan sains fisika, kimia dasar, dan biologi terintegrasi yang mudah dipahami melalui contoh kontekstual kehidupan sehari-hari.', 'https://images.unsplash.com/photo-1532094349884-543bc11b234d?w=600&auto=format&fit=crop&q=80', 'Kelas 5-8 SD/SMP', 'SMP', 180000.00, 125000.00, 1, 180, 1, NOW(), NOW());

-- 4. JADWAL LIVE CLASS RESMI
INSERT INTO live_classes (id, title, description, tutor_name, tutor_role, tutor_avatar, subject, grade, scheduled_at, duration_minutes, zoom_url, is_live, is_active, created_at, updated_at) VALUES
('l1111111-1111-1111-1111-111111111111', 'Mastering Trik Penalaran Aljabar & Aritmatika Sosial', 'Pembahasan metode bernalar khas KPM dalam memecahkan soal aljabar dan pemodelan matematika cerita.', 'Dr. Ridwan Hasan Saputra, M.Si', 'Founder & Master Coach KPM', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&auto=format&fit=crop&q=80', 'Matematika Nalaria', 'SMP & SMA', DATE_ADD(NOW(), INTERVAL 3 DAY), 90, 'https://zoom.us/j/8921471029?pwd=kpmliveclass', 0, 1, NOW(), NOW()),
('l2222222-2222-2222-2222-222222222222', 'Strategi Emas Juara Olimpiade Sains & MIPA Nasional', 'Bedah pola soal OSN/KSN tingkat kota hingga nasional beserta strategi manajemen waktu pengerjaan ujian CBT.', 'Kak Aris Fadillah, S.Si', 'Koordinator Pembinaan Olimpiade KPM', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&auto=format&fit=crop&q=80', 'Olimpiade MIPA', 'SD & SMP', DATE_ADD(NOW(), INTERVAL 4 DAY), 75, 'https://zoom.us/j/8921471030?pwd=kpmliveclass', 0, 1, NOW(), NOW()),
('l3333333-3333-3333-3333-333333333333', 'Trik Sukses Menjawab Soal Penalaran Matematika UTBK', 'Latihan bedah soal kuantitatif sulit dengan metode eliminasi logika dan penalaran praktis.', 'Kak Faiz Ramadhan, S.Pd', 'Tutor Spesialis Penalaran Skolastik', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&auto=format&fit=crop&q=80', 'Penalaran Kuantitatif', 'SMA / Alumni', DATE_ADD(NOW(), INTERVAL 5 DAY), 90, 'https://zoom.us/j/8921471031?pwd=kpmliveclass', 0, 1, NOW(), NOW());

-- 5. VIDEO MATERI PEMBELAJARAN RESMI
INSERT INTO videos (id, title, description, video_url, thumbnail, price, is_active, created_at, updated_at) VALUES
('v1111111-1111-1111-1111-111111111111', 'Pengenalan Konsep Dasar Berpikir Nalaria Realistik (MNR)', 'Penjelasan filosofi berpikir nalar, mengenali pola masalah matematika, dan cara berpikir sistematis.', 'https://www.youtube.com/watch?v=kJQP7kiw5Fk', 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=600&auto=format&fit=crop&q=80', 0.00, 1, NOW(), NOW()),
('v2222222-2222-2222-2222-222222222222', 'Trik Berhitung Cepat Operasi Bilangan Pecahan & Desimal', 'Metode penyelesaian pecahan campuran, perbandingan senilai, dan konversi persen secara cepat tanpa rumus rumit.', 'https://www.youtube.com/watch?v=kJQP7kiw5Fk', 'https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=600&auto=format&fit=crop&q=80', 0.00, 1, NOW(), NOW()),
('v3333333-3333-3333-3333-333333333333', 'Strategi Menyelesaikan Soal Cerita Persamaan Linear (SPLDV)', 'Langkah-langkah mengubah narasi soal cerita kehidupan sehari-hari menjadi bentuk persamaan matematika yang mudah dipecahkan.', 'https://www.youtube.com/watch?v=kJQP7kiw5Fk', 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=600&auto=format&fit=crop&q=80', 0.00, 1, NOW(), NOW()),
('v4444444-4444-4444-4444-444444444444', 'Kupas Tuntas Soal HOTS Logika Kombinatorika & Peluang', 'Analisis mendalam konsep permutasi, kombinasi, dan kaidah pencacahan untuk persiapan kompetisi sains dan ujian masuk PTN.', 'https://www.youtube.com/watch?v=kJQP7kiw5Fk', 'https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=600&auto=format&fit=crop&q=80', 0.00, 1, NOW(), NOW());

-- 6. MODUL & E-BOOK PDF RESMI
INSERT INTO modules (id, title, description, subject, grade, chapter, total_pages, pdf_url, cover_image, is_active, sort_order, created_at, updated_at) VALUES
('m1111111-1111-1111-1111-111111111111', 'Buku Panduan Konsep Matematika Nalaria Realistik (MNR) Dasar', 'Modul pengantar metode penalaran matematika realistik untuk tingkat SD dan SMP awal.', 'Matematika Nalaria', 'SD & SMP', 1, 36, 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf', 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=600&auto=format&fit=crop&q=80', 1, 1, NOW(), NOW()),
('m2222222-2222-2222-2222-222222222222', 'Kompilasi Soal Kompetisi Matematika Nalaria (KMNR) & Pembahasan', 'Kumpulan bank soal kompetisi KMNR babak penyisihan hingga final lengkap dengan pembahasan cara bernalar.', 'Olimpiade MIPA', 'Semua Jenjang', 2, 64, 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf', 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=600&auto=format&fit=crop&q=80', 1, 2, NOW(), NOW()),
('m3333333-3333-3333-3333-333333333333', 'Rangkuman Rumus Praktis & Pemecahan Soal MIPA SMP-SMA', 'Panduan saku rumus esensial matematika dan sains terpadu yang dilengkapi contoh soal latihan terarah.', 'Sains Terpadu', 'SMP & SMA', 3, 48, 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf', 'https://images.unsplash.com/photo-1532094349884-543bc11b234d?w=600&auto=format&fit=crop&q=80', 1, 3, NOW(), NOW()),
('m4444444-4444-4444-4444-444444444444', 'Bank Latihan Soal Tryout Penalaran Kuantitatif & TPS', 'Paket modul simulasi tes potensi skolastik dan penalaran numerik persiapan ujian seleksi perguruan tinggi.', 'Penalaran Kuantitatif', 'SMA / UTBK', 4, 52, 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf', 'https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=600&auto=format&fit=crop&q=80', 1, 4, NOW(), NOW());

-- 7. KUPON VOUCHER RESMI
INSERT INTO vouchers (id, code, title, description, discount_percent, discount_amount, min_purchase, max_discount, valid_until, is_active, usage_limit, usage_count, created_at, updated_at) VALUES
('c1111111-1111-1111-1111-111111111111', 'KPMJUARA2026', 'Diskon Spesial 25% Seluruh Paket', 'Potongan harga 25% untuk pembelian semua paket bimbingan dan tryout ujian CBT KPM Academy.', 25, 0.00, 0.00, 50000.00, DATE_ADD(NOW(), INTERVAL 1 YEAR), 1, 1000, 0, NOW(), NOW()),
('c2222222-2222-2222-2222-222222222222', 'MNRINDONESIA', 'Potongan Langsung Rp 50.000 Paket MNR', 'Spesial paket pembinaan Matematika Nalaria Realistik tingkat SD & SMP.', 0, 50000.00, 100000.00, 50000.00, DATE_ADD(NOW(), INTERVAL 1 YEAR), 1, 500, 0, NOW(), NOW()),
('c3333333-3333-3333-3333-333333333333', 'PELAJARHEBAT', 'Diskon 15% Member Baru KPM', 'Voucher potongan 15% untuk seluruh siswa baru KPM Academy.', 15, 0.00, 0.00, 35000.00, DATE_ADD(NOW(), INTERVAL 1 YEAR), 1, 500, 0, NOW(), NOW());
