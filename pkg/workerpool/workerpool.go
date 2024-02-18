package workerpool

// Pool Пул для выполнения
type Pool struct {
	// из этого канала будем брать выражения
	expressions chan map[string]interface{}
}

// New при создании пула передадим максимальное количество горутин и функция-обработчик
func New(maxGoroutines int, worker func(map[string]interface{})) *Pool {
	p := Pool{
		expressions: make(chan map[string]interface{}), // канал, откуда брать выражения
	}
	for i := 0; i < maxGoroutines; i++ {
		// создадим горутины по указанному количеству maxGoroutines
		go func() {
			// забираем выражения из канала
			for expression := range p.expressions {
				// и считаем
				worker(expression)
			}
		}()
	}

	return &p
}

func (p *Pool) Add(expression map[string]interface{}) {
	// добавляем выражения в канал, из которого забирает работу пул
	p.expressions <- expression
}

func (p *Pool) Shutdown() {
	// закроем канал с выражениями
	close(p.expressions)
}
