<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h3 my-4">Sepet</h1>
      </v-col>
    </v-row>

    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col cols="12" md="8">
        <v-card class="mb-4">
          <v-card-title class="text-subtitle-1 py-3">
            <div class="d-flex justify-space-between">
              <div>Ürün</div>
              <div>Toplam</div>
            </div>
          </v-card-title>
          <v-divider></v-divider>

          <v-card-text class="py-4" v-if="cartItems.length === 0">
            <p class="text-center">Sepetinizde ürün bulunmamaktadır.</p>
          </v-card-text>

          <v-card-text class="py-4" v-else>
            <div 
              v-for="(item, index) in cartItems" 
              :key="item.id"
              class="d-flex align-center"
              :class="{ 'mb-6': index < cartItems.length - 1 }"
            >
              <v-img
                :src="`https://picsum.photos/id/${100 + index}/100/100`"
                height="80"
                width="80"
                class="rounded me-4"
              ></v-img>
              <div class="flex-grow-1">
                <div class="text-h6">{{ item.name }}</div>
                <div class="text-body-2">₺{{ item.price }}</div>
                <div class="text-body-2">Stok: {{ item.stock }}</div>
                <div class="d-flex align-center mt-2">
                  <v-btn 
                    icon 
                    size="small" 
                    variant="text" 
                    density="comfortable"
                    @click="decreaseQuantity(item)"
                    :disabled="item.quantity <= 1"
                  >
                    <v-icon>mdi-minus</v-icon>
                  </v-btn>
                  <span class="mx-2">{{ item.quantity }}</span>
                  <v-btn 
                    icon 
                    size="small" 
                    variant="text" 
                    density="comfortable"
                    @click="increaseQuantity(item)"
                    :disabled="item.quantity >= item.stock"
                  >
                    <v-icon>mdi-plus</v-icon>
                  </v-btn>
                </div>
                <v-btn 
                  variant="text" 
                  density="comfortable" 
                  color="error" 
                  class="px-0 mt-2"
                  @click="removeItem(item)"
                >
                  <v-icon size="small" class="me-1">mdi-delete</v-icon>
                  Ürünü Kaldır
                </v-btn>
              </div>
              <div class="text-right">
                <div class="text-h6">₺{{ calculateItemTotal(item).toFixed(2) }}</div>
              </div>
            </div>
            <v-divider v-if="cartItems.length > 0"></v-divider>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="text-h5 py-4">Sepet Toplamı</v-card-title>
          
          <v-card-text>
            <v-expansion-panels variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title>Kupon Ekle</v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-text-field
                    label="Kupon Kodu"
                    variant="outlined"
                    density="comfortable"
                    hide-details
                    class="mb-2"
                    v-model="couponCode"
                  ></v-text-field>
                  <v-btn color="primary" block @click="applyCoupon">Uygula</v-btn>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>

            <v-list>
              <v-list-item>
                <div class="d-flex justify-space-between align-center py-2">
                  <div class="text-subtitle-1">Ara Toplam</div>
                  <div class="text-subtitle-1 font-weight-bold">₺{{ subtotal.toFixed(2) }}</div>
                </div>
              </v-list-item>
              <v-divider></v-divider>
              <v-list-item>
                <div class="d-flex justify-space-between align-center py-2">
                  <div class="text-h6 font-weight-bold">Toplam</div>
                  <div class="text-h6 font-weight-bold">₺{{ total.toFixed(2) }}</div>
                </div>
              </v-list-item>
            </v-list>

            <v-btn 
              color="success" 
              size="large" 
              block
              class="mt-4"
              :loading="processingOrder"
              :disabled="cartItems.length === 0 || processingOrder"
              @click="checkout"
            >
              Siparişi Tamamla
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';

interface Product {
  id: string;
  name: string;
  price: string;
  stock: number;
}

interface CartItem extends Product {
  quantity: number;
}

interface OrderItemRequest {
  product_id: string;
  quantity: number;
  price: string;
}

interface OrderRequest {
  user_id: string;
  order_items: OrderItemRequest[];
}

const loading = ref(true);
const cartItems = ref<CartItem[]>([]);
const couponCode = ref('');
const processingOrder = ref(false);

const fetchProducts = async () => {
  try {
    loading.value = true;
    const response = await axios.get('http://localhost:3001/api/product');
    
    // API'den gelen her ürünü sepete 1 adet olarak ekle
    cartItems.value = response.data.products.map((product: Product) => ({
      ...product,
      quantity: 1
    }));
  } catch (error) {
    console.error('Ürünler yüklenirken hata oluştu:', error);
  } finally {
    loading.value = false;
  }
};

const calculateItemTotal = (item: CartItem): number => {
  return parseFloat(item.price) * item.quantity;
};

const subtotal = computed((): number => {
  return cartItems.value.reduce((total, item) => {
    return total + calculateItemTotal(item);
  }, 0);
});

const total = computed((): number => {
  // Şimdilik ara toplam ile aynı, ileride kupon, vergi vs. eklenmesi durumunda değişecek
  return subtotal.value;
});

const increaseQuantity = (item: CartItem) => {
  if (item.quantity < item.stock) {
    item.quantity++;
  }
};

const decreaseQuantity = (item: CartItem) => {
  if (item.quantity > 1) {
    item.quantity--;
  }
};

const removeItem = (itemToRemove: CartItem) => {
  cartItems.value = cartItems.value.filter(item => item.id !== itemToRemove.id);
};

const applyCoupon = () => {
  // Kupon uygulama işlemi ileride API entegrasyonu ile eklenecek
  alert(`Kupon kodu uygulanıyor: ${couponCode.value}`);
  couponCode.value = '';
};

const checkout = async () => {
  if (cartItems.value.length === 0) {
    alert('Sepetinizde ürün bulunmamaktadır.');
    return;
  }

  try {
    processingOrder.value = true;

    // Sipariş verilerini API'nin beklediği formatta hazırla
    const orderItems: OrderItemRequest[] = cartItems.value.map(item => ({
      product_id: item.id,
      quantity: item.quantity,
      price: item.price
    }));

    const orderData: OrderRequest = {
      user_id: "test-user-id", // Gerçek uygulamada bu JWT token'dan alınacak
      order_items: orderItems
    };

    // Sipariş oluşturma isteğini gönder
    const response = await axios.post('http://localhost:3000/api/order/', orderData);
    
    // Başarılı sipariş durumunda sepeti temizle
    if (response.status === 200 || response.status === 201) {
      alert('Siparişiniz başarıyla oluşturuldu!');
      cartItems.value = [];
    }
  } catch (error: any) {
    console.error('Sipariş oluşturulurken hata oluştu:', error);
    alert(`Sipariş oluşturulamadı: ${error.response?.data?.message || 'Bilinmeyen hata'}`);
  } finally {
    processingOrder.value = false;
  }
};

onMounted(() => {
  fetchProducts();
});
</script>
