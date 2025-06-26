## чтобы создать флоу
Минимальный набор действий:
- добавить конфигурацию в файл
- зарегистрировать в методе registerFlows
```go
	flow := a.flowRegistry.CreateFlow("greeting")
```
- добавить условие старта флоу
```go
    flow.InitConditionFunc = func(c tele.Context) bool {
		return true
	}
```

## Опционально
- добавить колбэки на стейты
```go
flow.SetStateCallback("ask_phone", a.greetFlowCompletedCallback)
```
Колбэк будет вызван при получении инпута от пользователя на шаге