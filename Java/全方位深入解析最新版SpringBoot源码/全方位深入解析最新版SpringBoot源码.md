[TOC]



## 全局流程解析

### 万事俱备：SpringBoot项目环境准备

### SB的优越感：SpringBoot工程搭建演示

### 框架整体启动流程

框架初始化

框架启动

自动化装配

## 初始化器解析

### 系统初始化器实战

#### 实现方式一

##### 实现ApplicationContextInitializer接口

```java
@Order(1)
public class FirstInitializer implements ApplicationContextInitializer<ConfigurableApplicationContext> {
    @Override
    public void initialize(ConfigurableApplicationContext applicationContext) {
        ConfigurableEnvironment environment = applicationContext.getEnvironment();
//        environment.setRequiredProperties("mooc");
        Map<String, Object> map = new HashMap<>();
        map.put("key1", "value1");
        MapPropertySource mapPropertySource = new MapPropertySource("firstInitializer", map);
        environment.getPropertySources().addLast(mapPropertySource);
        System.out.println("run firstInitializer");
    }
}
```

##### 实现应用上下文接口

```java
@Component
public class TestService implements ApplicationContextAware {

    private ApplicationContext applicationContext;

    @Override
    public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
        this.applicationContext = applicationContext;
    }

    public String test() {
        return applicationContext.getEnvironment().getProperty("key3");
    }

}
```

##### spring.factories中配置接口实现

```properties
org.springframework.context.ApplicationContextInitializer=com.mooc.sb2.initializer.FirstInitializer
```

#### 实现方式二

##### SpringApplication类初始化时设置实现

```java
    public static void main(String[] args) {
//    SpringApplication.run(Sb2Application.class, args);
      SpringApplication springApplication = new SpringApplication(Sb2Application.class);
      springApplication.addInitializers(new SecondInitializer());
      springApplication.run(args);
   }
```

#### 实现方式三

##### application.properties中配置实现

```properties
context.initializer.classes=com.mooc.sb2.initializer.ThirdInitializer
```

### 工厂加载机制解析

#### SpringFactoriesLoader介绍

### 系统初始化器解析

## 监听器解析

### 实战监听器设计模式

### SpringBoot监听器实现

##### 