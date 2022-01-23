## Altastore

API online store sistem pickup single vendor

# User Endpoint :

```
http://localhost:1326/register
```

`POST` Endpoint untuk melakukan pendaftaran user, dengan request body seperti contoh dibawah.

| body :               |
|----------------------|
| "name" : "user1",    |
| "email" : "user1",   |
| "password" : "user1" |

```
http://localhost:1326/login
```

`POST` Endpoint untuk melakukan otentikasi pada pengguna yang telah mendaftar, untuk mendapatkan izin dan akses pada fitur dari API ini.

| body :               |
|----------------------|
| "email" : "user1",   |
| "password" : "user1" |

```
http://localhost:1326/users
```
`GET` Endpoint yang digunakan oleh Admin untuk melihat daftar user.

```
http://localhost:1326/users/profile
```

`GET` Pada Endpoint ini user dapat melihat data dirinya, pada endpoint ini.

```
http://localhost:1326/users/delete
```

`DELETE` Endpoint ini dapat digunakan oleh user untuk melakukan hapus akun, dengan melakukan otentikasi password terlebih dahulu.

| body :               |
|----------------------|
| "password" : "user1" |

```

http://localhost:1326/users/update
```

`PUT` User dapat melakukan perubahan terhadap datanya pada endpoint ini.

| body :                                 |
|----------------------------------------|
| "name" : "user1update", (Opsional)     |
| "email" : "user1update", (Opsional)    |
| "password" : "user1update", (Opsional) |

# Endpoint Category :
Endpoint category hanya bisa diakses oleh admin.

```
http://localhost:1326/categories
```


`GET` Method get pada endpoin ini digunakan untuk melihat daftar kategori.

`POST` Method untuk melakukan pembuatan kategori.

| body :                |
|-----------------------|
| "name" : "Handphone", |


```
http://localhost:8080/todos/:id
```

`GET` Method untuk menampilkan satu kategori berdasarkan ID.

`PUT` Method untuk mengedit Kategori.

| body :             |
|--------------------|
| "name" : "Laptop", |

`DELETE` Method untuk menghapus Kategori.


# Endpoint Product :


```
http://localhost:1326/product
```

`GET` Method get pada endpoin ini digunakan untuk melihat daftar produk, pada endpoint ini dapat ditambahkan query parameter untuk melakukan pencarian dan pagination. Query parameter tersebut bersifat opsional user dapat mengosongkannya dan sistem akan menaruh nilai default.

| query Parameter :                                  |
|----------------------------------------------------|
| "page=1", (Page 1 merupakan nilai default)         |
| "perpage=10", (PerPage 10 merupakan nilai default) |
| "search=", (Search memiliki nilai default NULL)    |

`POST` Method untuk melakukan pembuatan product.

| body :                       |
|------------------------------|
| "name" : "Xiaomi",           |
| "price" : 1000000,           |
| "stock" : 10,                |
| "Description" : "Hp Xiaomi", |
| "categoryID" : 1,            |


```
http://localhost:8080/todos/:id
```

`GET` Method untuk menampilkan satu Produk berdasarkan ID.

`PUT` Method untuk mengedit Produk.

| body :                                          |
|-------------------------------------------------|
| "name" : "Xiaomi 2", (OPSIONAL)                 |
| "price" : 1500000,  (OPSIONAL)                  |
| "stock" : 15,      (OPSIONAL)                   |
| "Description" : "Hp Xiaomi terbaru", (OPSIONAL) |
| "categoryID" : 1,(OPSIONAL)                     |

`DELETE` Method untuk menghapus Product.


# Endpoint Transaction :

```
http://localhost:1326/trnasaction
```

`GET` Method untuk menampilkan semua Transaksi. Jika user yang melakukan request maka yang tampil hanya transaksi yang dimiliki oleh user tersebut, tetapi jika admin yang melakukan request maka akan tampil semua transaksi dari semua user.

```
http://localhost:1326/trnasaction/checkout
```
`POST` Endpoint untuk melakukan pembayaran transaksi.

| body :            |
|-------------------|
| "product_id" : 1, |
| "quantity" : 5    |


```
http://localhost:1326/trnasaction/:transactionID
```
`GET` Method untuk menampilkan satu Transaksi berdasarkan transactionID yang dikirim melalui parameter.


# Endpoint Cart :

```
http://localhost:1326/cart
```
`GET` Pada endpoint ini user dapat melihat semua barang di keranjang belanja (Cart).

```
http://localhost:1326/cart/:productId
```
`POST` Dengan menggunakan method post user dapat menambahkan barang kedalam keranjang belanja (Cart).

| body :            |
|-------------------|
| "quantity" : 5    |

`PUT` User dapat merubah jumlah suatu barang pada keranjang belanja
| body :            |
|-------------------|
| "quantity" : 1    |

`DELETE` User dapat menghapus barang dari keranjang belanja menggunakan method DELETE pada endpoint ini.