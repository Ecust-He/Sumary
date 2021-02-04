[TOC]

## 函数式编程

### lambda表达式

#### 购物车案例

```java
/**
 * 购物车服务类
 */
public class CartService {

    // 加入到购物车中的商品信息
    private static List<Sku> cartSkuList =
            new ArrayList<Sku>(){
        {
            add(new Sku(654032, "无人机",
                    4999.00, 1,
                    4999.00, SkuCategoryEnum.ELECTRONICS));

            add(new Sku(642934, "VR一体机",
                    2299.00, 1,
                    2299.00, SkuCategoryEnum.ELECTRONICS));
        }
    };

    /**
     * 获取商品信息列表
     * @return
     */
    public static List<Sku> getCartSkuList() {
        return cartSkuList;
    }

    /**
     * Version 1.0.0
     * 找出购物车中所有电子产品
     * @param cartSkuList
     * @return
     */
    public static List<Sku> filterElectronicsSkus(
            List<Sku> cartSkuList) {

        List<Sku> result = new ArrayList<Sku>();
        for (Sku sku: cartSkuList) {
            // 如果商品类型 等于 电子类
            if (SkuCategoryEnum.ELECTRONICS.
                    equals(sku.getSkuCategory())) {
                result.add(sku);
            }
        }
        return result;
    }

    /**
     * Version 2.0.0
     * 根据传入商品类型参数，找出购物车中同种商品类型的商品列表
     * @param cartSkuList
     * @param category
     * @return
     */
    public static List<Sku> filterSkusByCategory(
            List<Sku> cartSkuList, SkuCategoryEnum category) {

        List<Sku> result = new ArrayList<Sku>();
        for (Sku sku: cartSkuList) {
            // 如果商品类型 等于 传入商品类型参数
            if (category.equals(sku.getSkuCategory())) {
                result.add(sku);
            }
        }
        return result;
    }

    /**
     * Version 3.0.0
     * 支持通过商品类型或总价来过滤商品
     * @param cartSkuList
     * @param category
     * @param totalPrice
     * @param categoryOrPrice - true: 根据商品类型，false: 根据商品总价
     * @return
     */
    public static List<Sku> filterSkus(
            List<Sku> cartSkuList, SkuCategoryEnum category,
            Double totalPrice, Boolean categoryOrPrice) {

        List<Sku> result = new ArrayList<Sku>();
        for (Sku sku: cartSkuList) {

            // 如果根据商品类型判断，sku类型与输入类型比较
            // 如果根据商品总价判断，sku总价与输入总价比较
            if (
                    (categoryOrPrice
                            && category.equals(sku.getSkuCategory())
                    ||
                    (!categoryOrPrice
                            && sku.getTotalPrice() > totalPrice))) {
                result.add(sku);
            }
        }
        return result;
    }

    /**
     * Version 4.0.0
     * 根据不同的Sku判断标准，对Sku列表进行过滤
     * @param cartSkuList
     * @param predicate - 不同的Sku判断标准策略
     * @return
     */
    public static List<Sku> filterSkus(
            List<Sku> cartSkuList, SkuPredicate predicate) {

        List<Sku> result = new ArrayList<Sku>();
        for (Sku sku: cartSkuList) {
            // 根据不同的Sku判断标准策略，对Sku进行判断
            if (predicate.test(sku)) {
                result.add(sku);
            }
        }
        return result;
    }
}
```

##### 版本一

```java
@Test
public void filterElectronicsSkus() {
    List<Sku> cartSkuList = CartService.getCartSkuList();

    // 查找购物车中数码类商品
    List<Sku> result =
            CartService.filterElectronicsSkus(cartSkuList);
}
```

##### 版本二

```java
@Test
public void filterSkusByCategory() {
    List<Sku> cartSkuList = CartService.getCartSkuList();

    // 查找购物车中图书类商品集合
    List<Sku> result = CartService.filterSkusByCategory(
            cartSkuList, SkuCategoryEnum.BOOKS);
}
```

##### 版本三

```java
@Test
public void filterSkus() {
    List<Sku> cartSkuList = CartService.getCartSkuList();

    // 根据商品总价过滤超过2000元的商品列表
    List<Sku> result = CartService.filterSkus(
            cartSkuList, null,
            2000.00, false);
}
```

##### 版本四

```java
@Test
public void filterSkus() {
    List<Sku> cartSkuList = CartService.getCartSkuList();

    // 过滤商品总价大于2000的商品
    List<Sku> result = CartService.filterSkus(
            cartSkuList, new SkuTotalPricePredicate());
}
```

##### 版本五

```java
@Test
public void filterSkus() {
    List<Sku> cartSkuList = CartService.getCartSkuList();

    // 过滤商品单价大于1000的商品
    List<Sku> result = CartService.filterSkus(
            cartSkuList, new SkuPredicate() {
                @Override
                public boolean test(Sku sku) {
                    return sku.getSkuPrice() > 1000;
                }
            });
}
```

##### 版本六

```java
@Test
public void filterSkus() {
    List<Sku> cartSkuList = CartService.getCartSkuList();

    // 过滤商品单价大于1000的商品
    List<Sku> result = CartService.filterSkus(
            cartSkuList,
            (Sku sku) -> sku.getSkuPrice() > 1000);
}
```

### 函数式接口

#### 文件处理案例（自定义函数式接口实现）

##### 文件处理函数式接口

```java
/**
 * 文件处理函数式接口
 */
@FunctionalInterface
public interface FileConsumer {

    /**
     * 函数式接口抽象方法
     * @param fileContent - 文件内容字符串
     */
    void fileHandler(String fileContent);

}
```

##### 文件服务类

```java
/**
 * 文件服务类
 */
public class FileService {

    /**
     * 从过url获取本地文件内容，调用函数式接口处理
     * @param url
     * @param fileConsumer
     */
    public void fileHandle(String url, FileConsumer fileConsumer)
            throws IOException {

        // 创建文件读取流
        BufferedReader bufferedReader = new BufferedReader(
                new InputStreamReader(
                        new FileInputStream(url)));

        // 定义行变量和内容sb
        String line;
        StringBuilder stringBuilder = new StringBuilder();

        // 循环读取文件内容
        while ((line = bufferedReader.readLine()) != null) {
            stringBuilder.append(line + "\n");
        }

        // 调用函数式接口方法，将文件内容传递给lambda表达式，实现业务逻辑
        fileConsumer.fileHandler(stringBuilder.toString());
    }
}
```

##### 测试

```java
@Test
public void fileHandle() throws IOException {
    FileService fileService = new FileService();

    // TODO 此处替换为本地文件的地址全路径
    String filePath = "";

    // 通过lambda表达式，打印文件内容
    fileService.fileHandle(filePath,

            System.out::println);
}
```

### 方法引用

#### 文件处理案例（自定义函数式接口实现）

```java
public class MethodReferenceTest {

    static class Sku {

        private String skuName;
        private Integer skuPrice;
        public Integer getSkuPrice() {
            return this.skuPrice;
        }

        public static int staticComparePrice(Sku sku1, Sku sku2) {
            return sku1.getSkuPrice() - sku2.getSkuPrice();
        }

        public int instanceComparePrice(Sku sku) {
            return this.getSkuPrice() - sku.getSkuPrice();
        }
    }

    class PriceComparator {
        public int instanceComparePrice(Sku sku1, Sku sku2) {
            return sku1.getSkuPrice() - sku2.getSkuPrice();
        }
    }

    public void test() {
        List<Sku> skuList = new ArrayList();

        skuList.sort((sku1, sku2) ->
                sku1.getSkuPrice() - sku2.getSkuPrice());

        // 类名::静态方法名
        skuList.sort(Sku::staticComparePrice);
        // 展开
        skuList.sort((Sku sku1, Sku sku2) -> {
            return Sku.staticComparePrice(sku1, sku2);
        });

        PriceComparator priceComparator = new PriceComparator();
        // 对象::实例方法名
        skuList.sort(priceComparator::instanceComparePrice);
        // 展开
        skuList.sort((Sku sku1, Sku sku2) -> {
            return priceComparator.instanceComparePrice(sku1, sku2);
        });

        // 类名::实例方法名
        skuList.sort(Sku::instanceComparePrice);
        // 展开
        skuList.sort((Sku object, Sku sku) -> {
            return object.instanceComparePrice(sku);
        });

        // 构造方法
        Optional.ofNullable(skuList)
                .orElseGet(ArrayList::new);

    }

}
```



