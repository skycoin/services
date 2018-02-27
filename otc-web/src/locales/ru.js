export default {
  header: {
    navigation: {
      whitepapers: 'Документация',
      downloads: 'Загрузки',
      explorer: 'Обозреватель',
      blog: 'Блог',
      roadmap: 'План',
    },
  },
  footer: {
    getStarted: 'Начать',
    explore: 'Дополнительно',
    community: 'Сообщество',
    wallet: 'Получить Кошелёк',
    infographics: 'Инфографика',
    whitepapers: 'Документация',
    blockchain: 'Блокчейн Обозреватель',
    blog: 'Блог',
    twitter: 'Twitter',
    reddit: 'Reddit',
    github: 'Github',
    telegram: 'Telegram',
    slack: 'Slack',
    roadmap: 'План',
    skyMessenger: 'Sky-Messenger',
    cxPlayground: 'CX Playground',
    team: 'Команда',
    subscribe: 'Рассылка',
    market: 'Markets',
    bitcoinTalks: 'Bitcointalks ANN',
    instagram: 'Instagram',
    facebook: 'Facebook',
    discord: 'Discord',
  },
  distribution: {
    rate: 'Current OTC rate: {rate} SKY/BTC',
    inventory: 'Current inventory: {coins} SKY available',
    title: 'Skycoin OTC',
    heading: 'Skycoin OTC',
    headingEnded: 'The previous distribution event finished on',
    ended: `<p>Join the <a href="https://t.me/skycoin">Skycoin Telegram</a>
      or follow the
      <a href="https://twitter.com/skycoinproject">Skycoin Twitter</a>
      to learn when the next event begins.`,
    instructions: `<p>You can check the current market value for <a href="https://coinmarketcap.com/currencies/skycoin/">Skycoin at CoinMarketCap</a>.</p>

<p>Что необходимо для участия в распространении:</p>

<ul>
  <li>Введите ваш Skycoin адрес</li>
  <li>Вы получите уникальный Bitcoin адрес для приобретения SKY</li>
  <li>Пошлите Bitcoin на полученый адрес</li>
</ul>

<p>Вы можете проверить статус заказа, введя адрес SKY и нажав на <strong>Проверить статус</strong>.</p>
<p>Каждый раз при нажатии на <strong>Получить адрес</strong>, генерируется новый BTC адрес. Один адрес SKY может иметь не более 5 BTC-адресов.</p>
    `,
    statusFor: 'Статус по {skyAddress}',
    enterAddress: 'Введите адрес Skycoin',
    getAddress: 'Получить адрес',
    checkStatus: 'Проверить статус',
    loading: 'Загрузка...',
    btcAddress: 'BTC адрес',
    errors: {
      noSkyAddress: 'Пожалуйста введите ваш SKY адрес.',
      coinsSoldOut: 'Skycoin OTC is currently sold out, check back later.',
    },
    statuses: {
      waiting_deposit: '[tx-{id} {updated}] Ожидаем BTC депозит.',
      waiting_send: '[tx-{id} {updated}] BTC депозит подтверждён. Skycoin транзакция поставлена в очередь.',
      waiting_confirm: '[tx-{id} {updated}] Skycoin транзакция отправлена. Ожидаем подтверждение.',
      done: '[tx-{id} {updated}] Завершена. Проверьте ваш Skycoin кошелёк.',
    },
  },
};
