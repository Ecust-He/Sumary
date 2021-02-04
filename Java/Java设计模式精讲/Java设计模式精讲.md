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