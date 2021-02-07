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

### 1、原始集合操作与Stream集合操作

```java
/**
 * 对比：原始集合操作与Stream集合操作
 */
public class StreamVs {

    /**
     * 1 想看看购物车中都有什么商品
     * 2 图书类商品都给买
     * 3 其余的商品中买两件最贵的
     * 4 只需要两件商品的名称和总价
     */

    /**
     * 以原始集合操作实现需求
     */
    @Test
    public void oldCartHandle() {
        List<Sku> cartSkuList = CartService.getCartSkuList();

        /**
         * 1 打印所有商品
         */
        for (Sku sku: cartSkuList) {
            System.out.println(JSON.toJSONString(sku, true));
        }

        /**
         * 2 图书类过滤掉
         */
        List<Sku> notBooksSkuList = new ArrayList<Sku>();
        for (Sku sku: cartSkuList) {
            if (!SkuCategoryEnum.BOOKS.equals(sku.getSkuCategory())) {
                notBooksSkuList.add(sku);
            }
        }

        /**
         * 排序
         */
        notBooksSkuList.sort(new Comparator<Sku>() {
            @Override
            public int compare(Sku sku1, Sku sku2) {
                if (sku1.getTotalPrice() > sku2.getTotalPrice()) {
                    return -1;
                } else if (sku1.getTotalPrice() < sku2.getTotalPrice()) {
                    return 1;
                } else {
                    return 0;
                }
            }
        });

        /**
         * TOP2
         */
        List<Sku> top2SkuList = new ArrayList<Sku>();
        for (int i = 0; i < 2; i++) {
            top2SkuList.add(notBooksSkuList.get(i));
        }

        /**
         * 4 求两件商品的总价
         */
        Double money = 0.0;
        for (Sku sku: top2SkuList) {
            // money = money + sku.getTotalPrice();
            money += sku.getTotalPrice();
        }

        /**
         * 获取两件商品的名称
         */
        List<String> resultSkuNameList = new ArrayList<String>();
        for (Sku sku: top2SkuList) {
            resultSkuNameList.add(sku.getSkuName());
        }

        /**
         * 打印输入结果
         */
        System.out.println(
                JSON.toJSONString(resultSkuNameList, true));
        System.out.println("商品总价：" + money);

    }

    /**
     * 以Stream流方式实现需求
     */
    @Test
    public void newCartHandle() {
        AtomicReference<Double> money =
                new AtomicReference<>(Double.valueOf(0.0));

        List<String> resultSkuNameList =
                CartService.getCartSkuList()
                .stream()
                /**
                 * 1 打印商品信息
                 */
                .peek(sku -> System.out.println(
                        JSON.toJSONString(sku, true)))
                /**
                 * 2 过滤掉所有图书类商品
                 */
                .filter(sku -> !SkuCategoryEnum.BOOKS.equals(
                        sku.getSkuCategory()))
                /**
                 * 排序
                 */
                .sorted(Comparator.
                        comparing(Sku::getTotalPrice).reversed())
                /**
                 * TOP2
                 */
                .limit(2)

                /**
                 * 累加商品总金额
                 */
                .peek(sku -> money.set(money.get() + sku.getTotalPrice()))

                /**
                 * 获取商品名称
                 */
                .map(sku -> sku.getSkuName())

                /**
                 * 收集结果
                 */
                .collect(Collectors.toList());


        /**
         * 打印输入结果
         */
        System.out.println(
                JSON.toJSONString(resultSkuNameList, true));
        System.out.println("商品总价：" + money.get());
    }

}
```

### 2、流的各种操作

```java
/**
 * 演示流的各种操作
 */
public class StreamOperator {

    List<Sku> list;

    @Before
    public void init() {
        list = CartService.getCartSkuList();
    }

    /**
     * filter使用：过滤掉不符合断言判断的数据
     */
    @Test
    public void filterTest() {
        list.stream()

                // filter
                .filter(sku ->
                        SkuCategoryEnum.BOOKS
                                .equals(sku.getSkuCategory()))

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * map使用：将一个元素转换成另一个元素
     */
    @Test
    public void mapTest() {
        list.stream()

                // map
                .map(sku -> sku.getSkuName())

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * flatMap使用：将一个对象转换成流
     */
    @Test
    public void flatMapTest() {
        list.stream()

                // flatMap
                .flatMap(sku -> Arrays.stream(
                        sku.getSkuName().split("")))

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * peek使用：对流中元素进行遍历操作，与forEach类似，但不会销毁流元素
     */
    @Test
    public void peek() {
        list.stream()

                // peek
                .peek(sku -> System.out.println(sku.getSkuName()))

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * sort使用：对流中元素进行排序，可选则自然排序或指定排序规则。有状态操作
     */
    @Test
    public void sortTest() {
        list.stream()

                .peek(sku -> System.out.println(sku.getSkuName()))

                //sort
                .sorted(Comparator.comparing(Sku::getTotalPrice))

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * distinct使用：对流元素进行去重。有状态操作
     */
    @Test
    public void distinctTest() {
        list.stream()
                .map(sku -> sku.getSkuCategory())

                // distinct
                .distinct()

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));


    }

    /**
     * skip使用：跳过前N条记录。有状态操作
     */
    @Test
    public void skipTest() {
        list.stream()

                .sorted(Comparator.comparing(Sku::getTotalPrice))

                // skip
                .skip(3)

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * limit使用：截断前N条记录。有状态操作
     */
    @Test
    public void limitTest() {
        list.stream()
                .sorted(Comparator.comparing(Sku::getTotalPrice))

                .skip(2 * 3)

                // limit
                .limit(3)

                .forEach(item ->
                        System.out.println(
                                JSON.toJSONString(
                                        item, true)));
    }

    /**
     * allMatch使用：终端操作，短路操作。所有元素匹配，返回true
     */
    @Test
    public void allMatchTest() {
        boolean match = list.stream()

                .peek(sku -> System.out.println(sku.getSkuName()))

                // allMatch
                .allMatch(sku -> sku.getTotalPrice() > 100);

        System.out.println(match);
    }

    /**
     * anyMatch使用：任何元素匹配，返回true
     */
    @Test
    public void anyMatchTest() {
        boolean match = list.stream()

                .peek(sku -> System.out.println(sku.getSkuName()))

                // anyMatch
                .anyMatch(sku -> sku.getTotalPrice() > 100);

        System.out.println(match);
    }

    /**
     * noneMatch使用：任何元素都不匹配，返回true
     */
    @Test
    public void noneMatchTest() {
        boolean match = list.stream()

                .peek(sku -> System.out.println(sku.getSkuName()))

                // noneMatch
                .noneMatch(sku -> sku.getTotalPrice() > 10_000);

        System.out.println(match);
    }

    /**
     * 找到第一个
     */
    @Test
    public void findFirstTest() {
        Optional<Sku> optional = list.stream()

                .peek(sku -> System.out.println(sku.getSkuName()))

                // findFirst
                .findFirst();

        System.out.println(
                JSON.toJSONString(optional.get(), true));
    }

    /**
     * 找任意一个
     */
    @Test
    public void findAnyTest() {
        Optional<Sku> optional = list.stream()

                .peek(sku -> System.out.println(sku.getSkuName()))

                // findAny
                .findAny();

        System.out.println(
                JSON.toJSONString(optional.get(), true));
    }

    /**
     * max使用：
     */
    @Test
    public void maxTest() {
        OptionalDouble optionalDouble = list.stream()
                // 获取总价
                .mapToDouble(Sku::getTotalPrice)

                .max();

        System.out.println(optionalDouble.getAsDouble());
    }

    /**
     * min使用
     */
    @Test
    public void minTest() {
        OptionalDouble optionalDouble = list.stream()
                // 获取总价
                .mapToDouble(Sku::getTotalPrice)

                .min();

        System.out.println(optionalDouble.getAsDouble());
    }

    /**
     * count使用
     */
    @Test
    public void countTest() {
        long count = list.stream()
                .count();

        System.out.println(count);
    }

}
```

### 3、流的四种构建形式

```java
/**
 * 流的四种构建形式
 */
public class StreamConstructor {

    /**
     * 由数值直接构建流
     */
    @Test
    public void streamFromValue() {
        Stream stream = Stream.of(1, 2, 3, 4, 5);

        stream.forEach(System.out::println);
    }

    /**
     * 通过数组构建流
     */
    @Test
    public void streamFromArray() {
        int[] numbers = {1, 2, 3, 4, 5};

        IntStream stream = Arrays.stream(numbers);
        stream.forEach(System.out::println);
    }

    /**
     * 通过文件生成流
     * @throws IOException
     */
    @Test
    public void streamFromFile() throws IOException {
        // TODO 此处替换为本地文件的地址全路径
        String filePath = "";

        Stream<String> stream = Files.lines(
                Paths.get(filePath));

        stream.forEach(System.out::println);
    }

    /**
     * 通过函数生成流（无限流）
     */
    @Test
    public void streamFromFunction() {

//        Stream stream = Stream.iterate(0, n -> n + 2);

        Stream stream = Stream.generate(Math::random);

        stream.limit(100)
                .forEach(System.out::println);

    }

}
```

### 4、常见收集器使用

```java
**
 * 常见预定义收集器使用
 */
public class StreamCollector {

    /**
     * 集合收集器
     */
    @Test
    public void toList() {

        List<Sku> list = CartService.getCartSkuList();

        List<Sku> result = list.stream()
                .filter(sku -> sku.getTotalPrice() > 100)

                .collect(Collectors.toList());

        System.out.println(
                JSON.toJSONString(result, true));

    }

    /**
     * 分组
     */
    @Test
    public void group() {
        List<Sku> list = CartService.getCartSkuList();

        // Map<分组条件，结果集合>
        Map<Object, List<Sku>> group = list.stream()
                .collect(
                        Collectors.groupingBy(
                                sku -> sku.getSkuCategory()));

        System.out.println(
                JSON.toJSONString(group, true));
    }

    /**
     * 分区
     */
    @Test
    public void partition() {
        List<Sku> list = CartService.getCartSkuList();

        Map<Boolean, List<Sku>> partition = list.stream()
                .collect(Collectors.partitioningBy(
                        sku -> sku.getTotalPrice() > 100));

        System.out.println(
                JSON.toJSONString(partition, true));
    }

}
```

### 5、规约和汇总操作

```java
public class ReduceAndCollectTest {


    @Test
    public void reduceTest() {

        /**
         * 订单对象
         */
        @Data
        @AllArgsConstructor
        class Order {
            /**
             * 订单编号
             */
            private Integer id;
            /**
             * 商品数量
             */
            private Integer productCount;
            /**
             * 消费总金额
             */
            private Double totalAmount;
        }

        /*
            准备数据
         */
        ArrayList<Order> list = Lists.newArrayList();
        list.add(new Order(1, 2, 25.12));
        list.add(new Order(2, 5, 257.23));
        list.add(new Order(3, 3, 23332.12));

        /*
            以前的方式：
            1. 计算商品数量
            2. 计算消费总金额
         */

        /*
            汇总商品数量和总金额
         */
        Order order = list.stream()
                .parallel()
                .reduce(
                        // 初始化值
                        new Order(0, 0, 0.0),

                        // Stream中两个元素的计算逻辑
                        (Order order1, Order order2) -> {
                            System.out.println("执行 计算逻辑 方法！！！");

                            int productCount =
                                    order1.getProductCount()
                                            + order2.getProductCount();

                            double totalAmount =
                                    order1.getTotalAmount()
                                            + order2.getTotalAmount();

                            return new Order(0, productCount, totalAmount);
                        },

                        // 并行情况下，多个并行结果如何合并
                        (Order order1, Order order2) -> {
                            System.out.println("执行 合并 方法！！！");

                            int productCount =
                                    order1.getProductCount()
                                            + order2.getProductCount();

                            double totalAmount =
                                    order1.getTotalAmount()
                                            + order2.getTotalAmount();

                            return new Order(0, productCount, totalAmount);
                        });

        System.out.println(JSON.toJSONString(order, true));
    }


    @Test
    public void collectTest() {
        /**
         * 订单对象
         */
        @Data
        @AllArgsConstructor
        class Order {
            /**
             * 订单编号
             */
            private Integer id;
            /**
             * 用户账号
             */
            private String account;
            /**
             * 商品数量
             */
            private Integer productCount;
            /**
             * 消费总金额
             */
            private Double totalAmount;
        }

        /*
            准备数据
         */
        ArrayList<Order> list = Lists.newArrayList();
        list.add(new Order(1, "zhangxiaoxi", 2, 25.12));
        list.add(new Order(2, "zhangxiaoxi",5, 257.23));
        list.add(new Order(3, "lisi",3, 23332.12));

        /*
            Map<用户账号, 订单(数量和金额)>
         */

        Map<String, Order> collect = list.stream()
                .parallel()
                .collect(
                        () -> {
                            System.out.println("执行 初始化容器 操作！！！");

                            return new HashMap<String, Order>();
                        },
                        (HashMap<String, Order> map, Order newOrder) -> {
                            System.out.println("执行 新元素添加到容器 操作！！！");

                            /*
                                新元素的account已经在map中存在了
                                不存在
                             */
                            String account = newOrder.getAccount();

                            // 如果此账号已存在，将新订单数据累加上
                            if (map.containsKey(account)) {
                                Order order = map.get(account);
                                order.setProductCount(
                                        newOrder.getProductCount()
                                                + order.getProductCount());
                                order.setTotalAmount(
                                        newOrder.getTotalAmount()
                                                + order.getTotalAmount());
                            } else {
                                // 如果不存在，直接将新订单存入map
                                map.put(account, newOrder);
                            }

                        }, (HashMap<String, Order> map1, HashMap<String, Order> map2) -> {
                            System.out.println("执行 并行结果合并 操作！！！");

                            map2.forEach((key, value) -> {
                                map1.merge(key, value, (order1, order2) -> {

                                    // TODO 注意：一定要用map1做合并，因为最后collect返回的是map1
                                    return new Order(0, key,
                                            order1.getProductCount()
                                                    + order2.getProductCount(),
                                            order1.getTotalAmount()
                                                    + order2.getTotalAmount());
                                });
                            });
                        });

        System.out.println(JSON.toJSONString(collect, true));
    }
}
```

### 6、实战案例（lambda+stream处理业务)

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

## 资源关闭

### 1、资源关闭优化前与优化后代码量对比

```java
/**
 * 资源关闭优化前与优化后代码量对比
 */
public class ResourceCloseVs {

    @Test
    public void newFileHandle(String url,
                              FileConsumer fileConsumer) {
        try (
                // 声明、创建文件的读取流
                FileInputStream fileInputStream =
                        new FileInputStream(url);

                InputStreamReader inputStreamReader =
                        new InputStreamReader(fileInputStream);

                BufferedReader bufferedReader =
                        new BufferedReader(inputStreamReader);
        ) {

            // 定义行变量和内容sb
            String line;
            StringBuilder stringBuilder = new StringBuilder();

            // 循环读取文件内容
            while ((line = bufferedReader.readLine()) != null) {
                stringBuilder.append(line + "\n");
            }

            // 调用函数式接口方法，将文件内容传递给lambda表达式，实现业务逻辑
            fileConsumer.fileHandler(stringBuilder.toString());

        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    @Test
    public void oldFileHandle(String url,
                              FileConsumer fileConsumer) {
        // 声明
        FileInputStream fileInputStream = null;
        InputStreamReader inputStreamReader = null;
        BufferedReader bufferedReader = null;

        // 创建文件读取流
        try {
            fileInputStream = new FileInputStream(url);

            inputStreamReader =
                    new InputStreamReader(fileInputStream);

            bufferedReader =
                    new BufferedReader(inputStreamReader);

            // 定义行变量和内容sb
            String line;
            StringBuilder stringBuilder = new StringBuilder();

            // 循环读取文件内容
            while ((line = bufferedReader.readLine()) != null) {
                stringBuilder.append(line + "\n");
            }

            // 调用函数式接口方法，将文件内容传递给lambda表达式，实现业务逻辑
            fileConsumer.fileHandler(stringBuilder.toString());

        } catch (IOException e) {
            e.printStackTrace();
        } finally {

            // 关闭流资源
            if (bufferedReader != null) {
                try {
                    bufferedReader.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
            if (inputStreamReader != null) {
                try {
                    inputStreamReader.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
            if (fileInputStream != null) {
                try {
                    fileInputStream.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
    }
}
```

## Guava工具

### 1、Optional用法

```java
/**
 * 学习Java8中的Optional使用方法
 */
public class OptionalTest {

    @Test
    public void test() throws Throwable {
        /**
         * 三种创建Optional对象方式
         */

        // 创建空的Optional对象
        Optional.empty();

        // 使用非null值创建Optional对象
        Optional.of("zhangxiaoxi");

        // 使用任意值创建Optional对象
        Optional optional = Optional.ofNullable("zhangxiaoxi");

        /**
         * 判断是否引用缺失的方法(建议不直接使用)
         */
        optional.isPresent();

        /**
         * 当optional引用存在时执行
         * 类似的方法：map filter flatMap
         */
        optional.ifPresent(System.out::println);


        /**
         * 当optional引用缺失时执行
         */
        optional.orElse("引用缺失");
        optional.orElseGet(() -> {
            // 自定义引用缺失时的返回值
            return "自定义引用缺失";
        });
        optional.orElseThrow(() -> {
            throw new RuntimeException("引用缺失异常");
        });
    }

    public static void stream(List<String> list) {
//        list.stream().forEach(System.out::println);


        Optional.ofNullable(list)
                .map(List::stream)
                .orElseGet(Stream::empty)
                .forEach(System.out::println);

    }

    public static void main(String[] args) {
        stream(Arrays.asList("1","2"));
    }
}
```

### 2、不可变集合用法

```java
/**
 * 不可变集合用法
 */
public class ImmutableTest {

    public static void test(List<Integer> list) {
        list.remove(0);
    }

    public static void main(String[] args) {
        List<Integer> list = new ArrayList<Integer>();

        list.add(1);
        list.add(2);
        list.add(3);

        List<Integer> newList =
                Collections.unmodifiableList(list);

        test(newList);

        System.out.println(newList);
    }

    public void immutable() {
        List<Integer> list = new ArrayList<Integer>();

        list.add(1);
        list.add(2);
        list.add(3);

        /**
         * 构造不可变集合对象三种方式
         */
        // 通过已经存在的集合创建
        ImmutableSet.copyOf(list);

        // 通过初始值，直接创建不可变集合
        ImmutableSet immutableSet =
                ImmutableSet.of(1, 2, 3);

        // 以builder方式创建
        ImmutableSet.builder()
                .add(1)
                .addAll(Sets.newHashSet(2, 3))
                .add(4)
                .build();
    }
}
```

### 3、新型集合Multiset用法

```java
/**
 * 实现：使用Multiset统计一首古诗的文字出现频率
 */
public class MultisetTest {

    private static final String text =
            "《南陵别儿童入京》" +
            "白酒新熟山中归，黄鸡啄黍秋正肥。" +
            "呼童烹鸡酌白酒，儿女嬉笑牵人衣。" +
            "高歌取醉欲自慰，起舞落日争光辉。" +
            "游说万乘苦不早，著鞭跨马涉远道。" +
            "会稽愚妇轻买臣，余亦辞家西入秦。" +
            "仰天大笑出门去，我辈岂是蓬蒿人。";

    @Test
    public void handle() {
        // multiset创建
        Multiset<Character> multiset =
                HashMultiset.create();

        // string 转换成 char 数组
        char[] chars = text.toCharArray();

        // 遍历数组，添加到multiset中
        Chars.asList(chars)
                .stream()
                .forEach(charItem -> {
                    multiset.add(charItem);
                });

        System.out.println
                ("size : " + multiset.size());

        System.out.println
                ("count : " + multiset.count('人'));
    }
}
```

### 4、Lists / Sets 使用

```java
/**
 * Lists / Sets 使用
 */
public class SetsTest {

    /**
     * Sets工具类的常用方法
     * 并集 / 交集 / 差集 / 分解集合中的所有子集 / 求两个集合的笛卡尔积
     *
     * Lists工具类的常用方式
     * 反转 / 拆分
     */

    private static final Set set1 =
            Sets.newHashSet(1, 2);

    private static final Set set2 =
            Sets.newHashSet(4);

    // 并集
    @Test
    public void union() {
        Set<Integer> set = Sets.union(set1, set2);

        System.out.println(set);
    }

    // 交集
    @Test
    public void intersection() {
        Set<Integer> set = Sets.intersection(set1, set2);

        System.out.println(set);
    }

    // 差集：如果元素属于A而且不属于B
    @Test
    public void difference() {
        Set<Integer> set = Sets.difference(set1, set2);

        System.out.println(set);

        // 相对差集：属于A而且不属于B 或者 属于B而且不属于A
        set = Sets.symmetricDifference(set1, set2);

        System.out.println(set);
    }

    // 拆分所有子集合
    @Test
    public void powerSet() {
        Set<Set<Integer>> powerSet = Sets.powerSet(set1);

        System.out.println(JSON.toJSONString(powerSet));
    }

    // 计算两个集合笛卡尔积
    @Test
    public void cartesianProduct() {
        Set<List<Integer>> product =
                Sets.cartesianProduct(set1, set2);

        System.out.println(JSON.toJSONString(product));
    }

    /**
     * 拆分
     */
    @Test
    public void partition() {
        List<Integer> list =
                Lists.newArrayList(1, 2, 3, 4, 5, 6, 7);

        List<List<Integer>> partition =
                Lists.partition(list, 3);

        System.out.println(JSON.toJSONString(partition));
    }

    // 反转
    @Test
    public void reverse() {
        List<Integer> list = Lists.newLinkedList();
        list.add(1);
        list.add(2);
        list.add(3);

        List<Integer> newList = Lists.reverse(list);

        System.out.println(newList);
    }

}
```

### 5、演示如何使用流(Source)与汇(Sink)来对文件进行常用操作

```java
/**
 * 演示如何使用流(Source)与汇(Sink)来对文件进行常用操作
 */
public class IOTest {
    @Test
    public void copyFile() throws IOException {
        /**
         * 创建对应的Source和Sink
         */
        CharSource charSource = Files.asCharSource(
                new File("SourceText.txt"),
                Charsets.UTF_8);
        CharSink charSink = Files.asCharSink(
                new File("TargetText.txt"),
                Charsets.UTF_8);
        /**
         * 拷贝
         */
        charSource.copyTo(charSink);
    }

}
```

### 6、演示布隆过滤器使用

```java
public class BloomFilterTest {


    @Test
    public void bloomFilter() {

        // 创建布隆过滤器
        BloomFilter<Integer> bloomFilter = BloomFilter.create(

                // 将任意类型数据转换为Java基础类型，默认转换为byte数组
                (Integer from, PrimitiveSink primitiveSink)
                        -> primitiveSink.putInt(from),

                // 预计插入的元素总数
                10000L,

                // 期望误判率(0.0 ~ 1.0)
                0.9
                );

        // 向布隆过滤器中添加元素
        for (int i = 0; i < 10000; i++) {
            bloomFilter.put(i);
        }

        // 检测给定元素是否 可能 存在在布隆过滤器中
        boolean might = bloomFilter.mightContain(6666);
        System.out.println("是否存在？" + might);
    }

}
```

## 线程池

### 1、线程池的简单使用

```java
    /**
     * 新的处理方式
     */
    @Test
    public void newHandle() throws InterruptedException {
        /**
         * 开启了一个线程池：线程个数是10个
         */
        ExecutorService threadPool =
                Executors.newFixedThreadPool(10);
        /**
         * 使用循环来模拟许多用户请求的场景
         */
        for (int request = 1; request <= 100; request++) {
            threadPool.execute(() -> {
                System.out.println("文档处理开始！");

                try {
                    // 将Word转换为PDF格式：处理时长很长的耗时过程
                    Thread.sleep(1000L * 30);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                System.out.println("文档处理结束！");
            });
        }
        Thread.sleep(1000L * 1000);
    }

    /**
     * 老的处理方式
     */
    @Test
    public void oldHandle() throws InterruptedException {
        /**
         * 使用循环来模拟许多用户请求的场景
         */
        for (int request = 1; request <= 100; request++) {
            new Thread(() -> {
                System.out.println("文档处理开始！");

                try {
                    // 将Word转换为PDF格式：处理时长很长的耗时过程
                    Thread.sleep(1000L * 30);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }

                System.out.println("文档处理结束！");
            }).start();
        }
        Thread.sleep(1000L * 1000);
    }
}
```

### 2、阻塞队列的选择

```java
public class QueueTest {

    @Test
    public void arrayBlockingQueue() throws InterruptedException {
        /**
         * 基于数组的有界阻塞队列，队列容量为10
         */
        ArrayBlockingQueue queue =
                new ArrayBlockingQueue<Integer>(10);

        // 循环向队列添加元素
        for (int i = 0; i < 20; i++) {
            queue.put(i);
            System.out.println("向队列中添加值：" + i);
        }
    }

    @Test
    public void linkedBlockingQueue() throws InterruptedException {
        /**
         * 基于链表的有界/无界阻塞队列，队列容量为10
         */
        LinkedBlockingQueue queue =
                new LinkedBlockingQueue<Integer>();

        // 循环向队列添加元素
        for (int i = 0; i < 20; i++) {
            queue.put(i);
            System.out.println("向队列中添加值：" + i);
        }
    }

    @Test
    public void test() throws InterruptedException {
        /**
         * 同步移交阻塞队列
         */
        SynchronousQueue queue = new SynchronousQueue<Integer>();

        // 插入值
        new Thread(() -> {
            try {
                queue.put(1);
                System.out.println("插入成功");
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }).start();

        // 删除值
        /*
        new Thread(() -> {
            try {
                queue.take();
                System.out.println("删除成功");
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }).start();
        */
        Thread.sleep(1000L * 60);
    }
}
```

### 3、饱和策略的选择

```java
public class PolicyTest {

    /**
     * 线程池
     */
    private static ThreadPoolExecutor executor =
            new ThreadPoolExecutor(
                    // 核心线程数和最大线程数
                    2, 3,

                    // 线程空闲后的存活时间
                    60L, TimeUnit.SECONDS,

                    // 有界阻塞队列
                    new LinkedBlockingQueue<Runnable>(5));

    /**
     * 任务
     */
    class Task implements Runnable {
        /**
         * 任务名称
         */
        private String taskName;

        public Task(String taskName) {
            this.taskName = taskName;
        }

        @Override
        public void run() {
            System.out.println("线程[ " + Thread.currentThread().getName()
                    + " ]正在执行[ " + this.taskName + " ]任务...");

            try {
                Thread.sleep(1000L * 5);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }

            System.out.println("线程[ " + Thread.currentThread().getName()
                    + " ]已执行完[ " + this.taskName + " ]任务！！！");
        }
    }

    /**
     * 线程池的执行过程
     *
     * 2个核心线程
     * 5个任务的队列
     * 3个最大线程：1个线程可用
     *
     * 前2个任务，会占用2个核心线程
     * 第3个到第7个任务，会暂存到任务队列中
     * 第8个任务，会启动最大线程，去执行
     * 第9个任务，没有线程可以去执行~~~
     */

    /**
     * 终止策略
     * TODO 抛出异常，拒绝任务提交
     */
    @Test
    public void abortPolicyTest() {
        // 设置饱和策略为 终止策略
        executor.setRejectedExecutionHandler(
                new ThreadPoolExecutor.AbortPolicy());

        for (int i = 1; i <= 10; i++) {
            try {
                // 提交10个线程任务
                executor.execute(new Task("线程任务" + i));
            } catch (Exception e) {
                System.err.println(e);
            }
        }

        // 关闭线程池
        executor.shutdown();
    }

    /**
     * 抛弃策略
     * TODO 直接丢弃掉新提交的任务
     */
    @Test
    public void discardPolicyTest() {
        // 设置饱和策略为 抛弃策略
        executor.setRejectedExecutionHandler(
                new ThreadPoolExecutor.DiscardPolicy());

        for (int i = 1; i <= 10; i++) {
            try {
                // 提交10个线程任务
                executor.execute(new Task("线程任务" + i));
            } catch (Exception e) {
                System.err.println(e);
            }
        }

        // 关闭线程池
        executor.shutdown();
    }

    /**
     * 抛弃旧任务策略
     * TODO 丢弃掉任务队列中的旧任务，暂存新提交的任务
     */
    @Test
    public void discardOldestPolicyTest() {
        // 设置饱和策略为 抛弃旧任务策略
        executor.setRejectedExecutionHandler(
                new ThreadPoolExecutor.DiscardOldestPolicy());

        for (int i = 1; i <= 10; i++) {
            try {
                // 提交10个线程任务
                executor.execute(new Task("线程任务" + i));
            } catch (Exception e) {
                System.err.println(e);
            }
        }

        // 关闭线程池
        executor.shutdown();
    }

    /**
     * 调用者运行策略
     * TODO 借用主线程来执行多余任务
     */
    @Test
    public void callerRunsPolicyTest() {
        // 设置饱和策略为 调用者运行策略
        executor.setRejectedExecutionHandler(
                new ThreadPoolExecutor.CallerRunsPolicy());

        for (int i = 1; i <= 10; i++) {
            try {
                // 提交10个线程任务
                executor.execute(new Task("线程任务" + i));
            } catch (Exception e) {
                System.err.println(e);
            }
        }

        // 关闭线程池
        executor.shutdown();
    }

    /**
     * 单元测试执行完，主线程等待100秒。防止主线程退出，看不到线程的执行结果
     * @throws InterruptedException
     */
    @After
    public void after() throws InterruptedException {
        Thread.sleep(1000L * 100);
    }
}
```

### 4、向线程池中提交任务

```java
public class RunTest {

    @Test
    public void submitTest()
            throws ExecutionException, InterruptedException {

        // 创建线程池
        ExecutorService threadPool =
                Executors.newCachedThreadPool();

        /**
         * 利用submit方法提交任务，接收任务的返回结果
         */
        Future<Integer> future = threadPool.submit(() -> {
            Thread.sleep(1000L * 10);

            return 2 * 5;
        });

        /**
         * 阻塞方法，直到任务有返回值后，才向下执行
         */
        Integer num = future.get();

        System.out.println("执行结果：" + num);
    }

    @Test
    public void executeTest() throws InterruptedException {
        // 创建线程池
        ExecutorService threadPool =
                Executors.newCachedThreadPool();

        /**
         * 利用execute方法提交任务，没有返回结果
         */
        threadPool.execute(() -> {
            try {
                Thread.sleep(1000L * 10);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }

            Integer num = 2 * 5;
            System.out.println("执行结果：" + num);
        });
        Thread.sleep(1000L * 1000);
    }
}
```