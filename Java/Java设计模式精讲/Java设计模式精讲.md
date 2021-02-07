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