[TOC]

## 七大设计原则

### 开闭原则

#### 编码演示

##### 课程接口

```java
public interface ICourse {
    Integer getId();
    String getName();
    Double getPrice();
}
```

##### Java课程

```java
public class JavaCourse implements ICourse{
    private Integer Id;
    private String name;
    private Double price;

    public JavaCourse(Integer id, String name, Double price) {
        this.Id = id;
        this.name = name;
        this.price = price;
    }

    public Double getPrice() {
        return this.price;
    }

}
```

##### Java打折课程

```java
public class JavaDiscountCourse extends JavaCourse {

    public JavaDiscountCourse(Integer id, String name, Double price) {
        super(id, name, price);
    }

    public Double getDiscountPrice(){
        return super.getPrice()*0.8;
    }
}
```

### 依赖倒置原则

#### 编码演示

##### 课程接口

```java
public interface ICourse {
    void studyCourse();
}
```

##### Java课程

```java
public class JavaCourse implements ICourse {

    @Override
    public void studyCourse() {
        System.out.println("Geely在学习Java课程");
    }
}
```

##### 学生类

```java
public class Geely {

    public void setiCourse(ICourse iCourse) {
        this.iCourse = iCourse;
    }

    private ICourse iCourse;

    public void studyImoocCourse(){
        iCourse.studyCourse();
    }
}
```

##### 测试类

```java
public class Test {

    //v1
//    public static void main(String[] args) {
//        Geely geely = new Geely();
//        geely.studyJavaCourse();
//        geely.studyFECourse();
//    }

    //v2
//    public static void main(String[] args) {
//        Geely geely = new Geely();
//        geely.studyImoocCourse(new JavaCourse());
//        geely.studyImoocCourse(new FECourse());
//        geely.studyImoocCourse(new PythonCourse());
//    }

    //v3
//    public static void main(String[] args) {
//        Geely geely = new Geely(new JavaCourse());
//        geely.studyImoocCourse();
//    }
    public static void main(String[] args) {
        Geely geely = new Geely();
        geely.setiCourse(new JavaCourse());
        geely.studyImoocCourse();

        geely.setiCourse(new FECourse());
        geely.studyImoocCourse();
    }
}
```

### 单一职责原则

#### 编码演示1

##### 鸟类

```java
public class Bird {
    public void mainMoveMode(String birdName){
        if("鸵鸟".equals(birdName)){
            System.out.println(birdName+"用脚走");
        }else{
            System.out.println(birdName+"用翅膀飞");
        }
    }
}
```

##### 会飞的鸟

```java
public class FlyBird {
    public void mainMoveMode(String birdName){
        System.out.println(birdName+"用翅膀飞");
    }
}
```

##### 会行走的鸟

```java
public class WalkBird {
    public void mainMoveMode(String birdName){
        System.out.println(birdName+"用脚走");
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
//        Bird bird = new Bird();
//        bird.mainMoveMode("大雁");
//        bird.mainMoveMode("鸵鸟");

        FlyBird flyBird = new FlyBird();
        flyBird.mainMoveMode("大雁");

        WalkBird walkBird = new WalkBird();
        walkBird.mainMoveMode("鸵鸟");

    }
}
```

#### 编码演示2

##### 课程接口类

```java
public interface ICourse {
    String getCourseName();
    byte[] getCourseVideo();

    void studyCourse();
    void refundCourse();

}
```

##### 课程内容接口类

```java
public interface ICourseContent {
    String getCourseName();
    byte[] getCourseVideo();
}
```

##### 课程管理接口类

```java
public interface ICourseManager {
    void studyCourse();
    void refundCourse();
}
```

##### 课程实现类

```java
public class CourseImpl implements ICourseManager,ICourseContent {
    @Override
    public void studyCourse() {

    }

    @Override
    public void refundCourse() {

    }

    @Override
    public String getCourseName() {
        return null;
    }

    @Override
    public byte[] getCourseVideo() {
        return new byte[0];
    }
}
```

### 接口隔离原则

#### 编码演示

```java
public interface IEatAnimalAction {
    void eat();
}
```

```java
public interface ISwimAnimalAction {
    void swim();
}
```

```java
public class Dog implements ISwimAnimalAction,IEatAnimalAction {

    @Override
    public void eat() {

    }

    @Override
    public void swim() {

    }
}
```

### 迪米特法则

#### 编码演示

##### 课程类

```java
public class Course {
}
```

##### 领导类

```java
public class TeamLeader {
    public void checkNumberOfCourses(){
        List<Course> courseList = new ArrayList<Course>();
        for(int i = 0 ;i < 20;i++){
            courseList.add(new Course());
        }
        System.out.println("在线课程的数量是："+courseList.size());
    }
}
```

##### 老板类

```java
public class Boss {

    public void commandCheckNumber(TeamLeader teamLeader){
        teamLeader.checkNumberOfCourses();
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Boss boss = new Boss();
        TeamLeader teamLeader = new TeamLeader();
        boss.commandCheckNumber(teamLeader);
    }
}
```

### 里式替换原则

#### 编码演示1

##### 父类

```java
public class Base {
    public void method(HashMap map){
        System.out.println("父类被执行");
    }
}
```

##### 子类

```java
public class Child extends Base {
//    @Override
//    public void method(HashMap map) {
//        System.out.println("子类HashMap入参方法被执行");
//    }

    public void method(Map map) {
        System.out.println("子类HashMap入参方法被执行");
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Base child = new Child();
        HashMap hashMap = new HashMap();
        child.method(hashMap);
    }
}
```

#### 编码演示2

##### 父类

```java
public abstract class Base {
    public abstract Map method();

}
```

##### 子类

```java
public class Child extends Base {
    @Override
    public HashMap method() {
        HashMap hashMap = new HashMap();
        System.out.println("子类method被执行");
        hashMap.put("message","子类method被执行");
        return hashMap;
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Child child = new Child();
        System.out.println(child.method());
    }
}
```

### 组合聚合原则

#### 编码演示

##### 数据库连接接口

```java
public abstract class DBConnection {
//    public String getConnection(){
//        return "MySQL数据库连接";
//    }
    public abstract String getConnection();
}
```

##### mysql数据库连接

```java
public class MySQLConnection extends DBConnection {
    @Override
    public String getConnection() {
        return "MySQL数据库连接";
    }
}
```

##### postgresql数据库连接

```java
public class PostgreSQLConnection extends DBConnection {
    @Override
    public String getConnection() {
        return "PostgreSQL数据库连接";
    }
}
```

##### 商品数据库操作

```java
public class ProductDao{
    private DBConnection dbConnection;

    public void setDbConnection(DBConnection dbConnection) {
        this.dbConnection = dbConnection;
    }

    public void addProduct(){
        String conn = dbConnection.getConnection();
        System.out.println("使用"+conn+"增加产品");
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        ProductDao productDao = new ProductDao();
        productDao.setDbConnection(new PostgreSQLConnection());
        productDao.addProduct();
    }
}
```