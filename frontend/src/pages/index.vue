<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h3 my-6 font-weight-light text-center text-md-start">Sepetim</h1>
      </v-col>
    </v-row>

    <v-row v-if="loading">
      <v-col cols="12" class="d-flex justify-center align-center py-12">
        <v-progress-circular indeterminate color="primary" size="60" width="6"></v-progress-circular>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col cols="12" md="8">
        <v-card class="mb-4 rounded-xl elevation-3 border" variant="elevated">
          <v-card-title class="text-subtitle-1 py-4 px-6 bg-grey-lighten-4 rounded-t-xl">
            <div class="d-flex justify-space-between">
              <div class="font-weight-medium">Ürün</div>
              <div class="font-weight-medium">Toplam</div>
            </div>
          </v-card-title>

          <v-card-text class="py-6 px-6" v-if="cartItems.length === 0">
            <div class="d-flex flex-column align-center py-8">
              <v-icon icon="mdi-cart-outline" size="64" color="grey-lighten-2" class="mb-4"></v-icon>
              <p class="text-body-1 text-medium-emphasis">Sepetinizde ürün bulunmamaktadır.</p>
            </div>
          </v-card-text>

          <v-card-text class="py-6 px-6" v-else>
            <div v-for="(item, index) in cartItems" :key="item.id">
              <div class="d-flex flex-column flex-md-row align-md-center rounded-lg hover-elevation transition-all">
                <div class="product-image-container d-flex justify-center justify-md-start mb-3 mb-md-0">
                  <div class="img-square-wrapper">
                    <v-img
                      src="/product.png"
                      width="90"
                      aspect-ratio="1"
                      class="rounded-lg bg-grey-lighten-4 object-cover border"
                      cover
                    ></v-img>
                  </div>
                </div>
                <div class="flex-grow-1 px-md-4">
                  <div class="text-h6 font-weight-medium mb-1 text-center text-md-start">{{ item.name }}</div>
                  <div class="text-subtitle-1 text-primary font-weight-medium mb-1 text-center text-md-start">₺{{ item.price }}</div>
                  <div class="text-caption text-medium-emphasis text-center text-md-start">Stok: {{ item.stock }}</div>
                  <div class="d-flex align-center justify-center justify-md-start mt-3">
                    <v-btn 
                      icon 
                      size="small" 
                      variant="outlined" 
                      color="grey-darken-1"
                      density="comfortable"
                      @click="decreaseQuantity(item)"
                      :disabled="item.quantity <= 1"
                      class="quantity-btn"
                    >
                      <v-icon>mdi-minus</v-icon>
                    </v-btn>
                    <span class="mx-4 text-subtitle-2 font-weight-medium">{{ item.quantity }}</span>
                    <v-btn 
                      icon 
                      size="small" 
                      variant="outlined" 
                      color="grey-darken-1"
                      density="comfortable"
                      @click="increaseQuantity(item)"
                      :disabled="item.quantity >= item.stock"
                      class="quantity-btn"
                    >
                      <v-icon>mdi-plus</v-icon>
                    </v-btn>
                  </div>
                </div>
                <div class="text-right d-flex flex-column align-center align-md-end mt-3 mt-md-0">
                  <div class="text-h6 font-weight-bold">₺{{ calculateItemTotal(item).toFixed(2) }}</div>
                  <v-btn 
                    variant="text" 
                    density="comfortable" 
                    color="error" 
                    class="mt-2 px-2"
                    @click="removeItem(item)"
                    size="small"
                  >
                    <v-icon size="small" class="me-1">mdi-delete</v-icon>
                    Kaldır
                  </v-btn>
                </div>
              </div>
              <v-divider v-if="index < cartItems.length - 1" class="my-4"></v-divider>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card class="rounded-xl elevation-3 border" variant="elevated">
          <v-card-title class="text-h6 py-4 px-6 bg-grey-lighten-4 rounded-t-xl font-weight-medium">
            Sipariş Özeti
          </v-card-title>
          
          <v-card-text class="px-6 py-4">
            <v-list class="pa-0">
              <v-list-item class="px-0">
                <div class="d-flex justify-space-between align-center py-2">
                  <div class="text-subtitle-2 text-medium-emphasis">Ara Toplam</div>
                  <div class="text-subtitle-1 font-weight-medium">₺{{ subtotal.toFixed(2) }}</div>
                </div>
              </v-list-item>
              <v-list-item class="px-0">
                <div class="d-flex justify-space-between align-center py-2">
                  <div class="text-subtitle-2 text-medium-emphasis">Kargo</div>
                  <div class="text-subtitle-1 font-weight-medium">Ücretsiz</div>
                </div>
              </v-list-item>
              <v-divider class="my-2"></v-divider>
              <v-list-item class="px-0">
                <div class="d-flex justify-space-between align-center py-3">
                  <div class="text-h6 font-weight-bold">Toplam</div>
                  <div class="text-h6 font-weight-bold text-primary">₺{{ total.toFixed(2) }}</div>
                </div>
              </v-list-item>
            </v-list>

            <v-btn 
              color="success" 
              size="large" 
              block
              class="mt-6 py-3 text-subtitle-1 rounded-lg"
              :loading="processingOrder"
              :disabled="cartItems.length === 0 || processingOrder"
              @click="checkout"
              elevation="1"
            >
              <v-icon class="me-2">mdi-check-circle</v-icon>
              Siparişi Tamamla
            </v-btn>
            
            <div class="d-flex align-center justify-center mt-4 gap-2">
              <v-icon size="small" color="grey-darken-1">mdi-shield-lock</v-icon>
              <span class="text-caption text-grey-darken-1">Güvenli Ödeme</span>
            </div>
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

<style scoped>
.hover-elevation {
  transition: all 0.2s ease-in-out;
  padding: 16px;
  border: 1px solid transparent;
}

.hover-elevation:hover {
  background-color: rgba(0, 0, 0, 0.01);
  border-color: rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.transition-all {
  transition: all 0.2s ease;
}

.border {
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.quantity-btn {
  transition: all 0.2s ease;
  min-width: 36px;
  min-height: 36px;
}

.quantity-btn:hover:not(:disabled) {
  background-color: rgba(0, 0, 0, 0.05);
}

.product-image-container {
  min-width: 90px;
  max-width: 90px;
}

.img-square-wrapper {
  width: 90px;
  height: 90px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

@media (max-width: 600px) {
  .hover-elevation {
    padding: 12px 8px;
  }
  
  .quantity-btn {
    min-width: 40px;
    min-height: 40px;
  }
  
  .product-image-container {
    min-width: 100%;
    max-width: 100%;
  }
  
  .img-square-wrapper {
    width: 90px;
    height: 90px;
    margin: 0 auto;
  }
}
</style>
