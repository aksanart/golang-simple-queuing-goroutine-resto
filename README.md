## Belajar simple queuing goroutine, study kasus antrian customer direstoran
### rule
1. restoran
    1. punya 5 meja (bisa menampung 5 orang)
    1. terbuka selama 30 (bisa dinamis) detik
    1. punya staff (waiter, chef)
    1. punya pelanggan
1. pelanggan
    1. pelanggan baru datang interval 3-5 (bisa dinamis) detik
    1. pelanggan duduk dimeja yang kosong, jika tidak ada meja yang kosong pelanggan pulang
    1. pelanggan memesan pesanan
    1. pelanggan menunggu pesanan 2 (bisa dinamis) detik
    1. pelanggan makan sekitar 3 (bisa dinamis) detik, setelah itu pulang
1. waiter
    1. waiter ada x (bisa dinamis) orang
    1. waiter mengecek pelanggan jika tidak ada istirahat jika ada catat pesanan
    1. waiter memberi Chef catatan pesanan
    1. waiter jika ada sinyal dari Chef maka waiter mengantarkan pesanan pelanggan
1. chef
    1. chef ada x (bisa dinamis) orang
    1. chef mengecek catatan pesanan jika tidak ada maka istirahat jika ada maka masak pesanan
    1. chef membunyikan bel sinyal ke waiter mengambil pesanan
