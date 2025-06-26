## чтобы создать флоу

Минимальный набор действий:

- добавить конфигурацию в файл
- зарегистрировать в методе registerFlows

```go
    flow := a.flowRegistry.CreateFlow("greeting")
```

- добавить условие старта флоу

```go
    flow.InitConditionFunc = func (c tele.Context) bool {
    return true
}
```

## Опционально

- добавить коллбеки на стейты

1. В routes.go добавить эндпоинт c кастомным названием

```go
a.bot.Handle("do_something", func (c tele.Context) error {
    a.handleFlow(c, nil)
    return c.Send("done")
})

```

> Важно: колбэк не дополняет дефолтный, а заменяет его.
> Поэтому при назначении кастомного обработчика нужно не забыть вызвать `a.handleFlow(c, nil)`, чтобы корректно обработать шаг флоу

2. в `config.yaml` написать action для нужного шага

```yaml
steps:
    first:
        action: do_something
        message:
            text: Answer the question
```