package spike

import "context"

// Spike 秒杀
type Spike interface {
	// 获取详情
	Detail() (Content, error)
	// 抢
	Grab(ctx context.Context) (uint, error)
}

type Content struct {
}

func NewSpike() Spike {
	return &spikeImpl{}
}

type spikeImpl struct {
}

// 获取详情
func (impl *spikeImpl) Detail() (Content, error) {
	// 从缓存获取

	return Content{}, nil
}

// 抢
func (impl *spikeImpl) Grab(ctx context.Context) (uint, error) {
	// 先获取锁
	{
		// 物品库存校验

		// 事务
		{
			// 物品库存减1

			// 新增订单--表明用户有购买权，在15分钟内支付完成，即可获得该物品；如果用户没有支付，则标记此订单为失效，还需将物品库存数回退
		}

		// 释放锁
	}

	// 返回订单id

	return 0, nil
}
