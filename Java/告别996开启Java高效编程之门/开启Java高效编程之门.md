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

##### 演示

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

#### 演示

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

## 流式编程

### 实战（lambda+stream处理业务)

#### 案例一：anyMatch的使用方式

```java
    @Test
    public void findStudent() {

        studentMap.forEach((studentName, scoreList) -> {

            boolean bool = scoreList.stream()
                    .anyMatch(score -> {
                        // TODO anyMatch找到任意一条符合条件的数据后就停止
//                        System.out.println(score);

                        return score.getScoreValue() == null;
                    });

            if (bool) {
                System.out.println("此学生[ " + studentName + " ]有缺考情况！");
            }
        });
    }
```

#### 案例二：filter和distinct的使用方式

```java
@Test
public void distinctTag() {

    tagListFromReq.stream()

            // TODO true:通过测试，数据不过滤；false:未通过测试，数据被过滤
            .filter(tag -> !tagListFromDB.contains(tag.getName()))

            // TODO 使用equals对元素进行比较
            .distinct()
            .forEach(tag -> System.out.println(tag));

}
```

#### 案例三：flatMap的使用方式

```java
    @Test
    public void findPermission() {
        roleList.stream()

                // TODO 扁平化MAP 获取对象中的集合类属性，组成一个新的流
                .flatMap(role -> role.getPermissions().stream())

                // peek 与 forEach 类似，区别是用在中间过程中，后面可以接其他操作
                .peek(permission ->
                        System.out.println("新的流元素：" + permission))

                .distinct()
//                .forEach(permission -> System.out.println(permission));
                .collect(Collectors.toList());
    }
```

#### 案例四：group的使用方式

```java
/**
 * 接口
 * @param accountIds
 * @return
 */
public Map<String, List<Order>> queryOrderByAccountIds(
        List<String> accountIds) {

    return Optional.ofNullable(selectFromDB(accountIds))
            .map(List::stream)
            .orElseGet(Stream::empty)

            // TODO group分组功能
            .collect(Collectors.groupingBy(
                    order -> order.getAccountId()));
}

@Test
public void test() {
    Map<String, List<Order>> orders =
            queryOrderByAccountIds(
                    Lists.newArrayList("张三", "李四", "王五"));

    System.out.println(JSON.toJSONString(orders, true));
}
```

#### 案例五：sort和compare的使用方式

```java
@Test
public void sortTrade() {

    System.out.println("排序前数据~~~\n" + JSON.toJSONString(trades, true));

    List<Trade> sorted = trades.stream()
            .sorted(
                    Comparator
                            // 首先按照价格排序
                            .comparing(
                                    Trade::getPrice,
                                    // TODO 进行排序调整，将自然排序翻转
                                    Comparator.reverseOrder())

                            // 时间先后进行排序，自然排序
                            .thenComparing(Trade::getTime)

                            // 交易量排序，自然排序翻转
                            .thenComparing(
                                    Trade::getCount,
                                    Comparator.reverseOrder())

                            // 自定义排序规则
                            .thenComparing(
                                    // 要排序的字段值
                                    Trade::getType,

                                    // 自定义排序规则
                                    (type1, type2) -> {
                                        if ("机构".equals(type1) && "个人".equals(type2)) {
                                            // -1:type1在先， type2在后
                                            return -1;
                                        } else if ("个人".equals(type1) && "机构".equals(type2)) {
                                            return 1;
                                        } else {
                                            return 0;
                                        }
                                    }))
            .collect(Collectors.toList());

    System.out.println("排序后结果~~~\n" + JSON.toJSONString(sorted, true));
}
```