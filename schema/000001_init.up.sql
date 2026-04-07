CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL DEFAULT uuidv7(),
    telegram_user_id BIGINT NOT NULL,
    telegram_user_name VARCHAR(32) DEFAULT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

COMMENT ON TABLE messages IS 'Таблица сообщений отправленных в telegram бот';
COMMENT ON COLUMN messages.id IS 'Идентификатор сообщения';
COMMENT ON COLUMN messages.telegram_user_id IS 'Идентификатор telegram пользователя';
COMMENT ON COLUMN messages.telegram_user_name IS 'Имя telegram пользователя';
COMMENT ON COLUMN messages.text IS 'Текст сообщения';
COMMENT ON COLUMN messages.created_at IS 'дата и время создания записи';

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'sunt', 'quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'molestias', 'et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eum', 'ullam et saepe reiciendis voluptatem adipisci\nsit amet autem assumenda provident rerum culpa\nquis hic commodi nesciunt rem tenetur doloremque ipsam iure\nquis sunt voluptatem rerum illo velit');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'nesciunt', 'repudiandae veniam quaerat sunt sed\nalias aut fugiat sit autem sed est\nvoluptatem omnis possimus esse voluptatibus quis\nest aut tenetur dolor neque');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolorem', 'ut aspernatur corporis harum nihil quis provident sequi\nmollitia nobis aliquid molestiae\nperspiciatis et ea nemo ab reprehenderit accusantium quas\nvoluptate dolores velit et doloremque molestiae');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'magnam', 'dolore placeat quibusdam ea quo vitae\nmagni quis enim qui quis quo nemo aut saepe\nquidem repellat excepturi ut quia\nsunt ut sequi eos ea sed quas');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolorem', 'dignissimos aperiam dolorem qui eum\nfacilis quibusdam animi sint suscipit qui sint possimus cum\nquaerat magni maiores excepturi\nipsam ut commodi dolor voluptatum modi aut vitae');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'nesciunt', 'consectetur animi nesciunt iure dolore\nenim quia ad\nveniam autem ut quam aut nobis\net est aut quod aut provident voluptas autem voluptas');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'optio', 'quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'etea', 'delectus reiciendis molestiae occaecati non minima eveniet qui voluptatibus\naccusamus in eum beatae sit\nvel qui neque voluptates ut commodi qui incidunt\nut animi commodi');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'quibusdam', 'itaque id aut magnam\npraesentium quia et ea odit et ea voluptas et\nsapiente quia nihil amet occaecati quia id voluptatem\nincidunt ea est distinctio odio');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolorum', 'aut dicta possimus sint mollitia voluptas commodi quo doloremque\niste corrupti reiciendis voluptatem eius rerum\nsit cumque quod eligendi laborum minima\nperferendis recusandae assumenda consectetur porro architecto ipsum ipsam');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'voluptatem', 'fuga et accusamus dolorum perferendis illo voluptas\nnon doloremque neque facere\nad qui dolorum molestiae beatae\nsed aut voluptas totam sit illum');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eveniet', 'reprehenderit quos placeat\nvelit minima officia dolores impedit repudiandae molestiae nam\nvoluptas recusandae quis delectus\nofficiis harum fugiat vitae');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'sint', 'suscipit nam nisi quo aperiam aut\nasperiores eos fugit maiores voluptatibus quia\nvoluptatem quis ullam qui in alias quia est\nconsequatur magni mollitia accusamus ea nisi voluptate dicta');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'fugit', 'eos voluptas et aut odit natus earum\naspernatur fuga molestiae ullam\ndeserunt ratione qui eos\nqui nihil ratione nemo velit ut aut id quo');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'voluptate', 'eveniet quo quis\nlaborum totam consequatur non dolor\nut et est repudiandae\nest voluptatem vel debitis et magnam');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'adipisci', 'illum quis cupiditate provident sit magnam\nea sed aut omnis\nveniam maiores ullam consequatur atque\nadipisci quo iste expedita sit quos voluptas');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'doloribus', 'qui consequuntur ducimus possimus quisquam amet similique\nsuscipit porro ipsam amet\neos veritatis officiis exercitationem vel fugit aut necessitatibus totam\nomnis rerum consequatur expedita quidem cumque explicabo');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'asperiores', 'repellat aliquid praesentium dolorem quo\nsed totam minus non itaque\nnihil labore molestiae sunt dolor eveniet hic recusandae veniam\ntempora et tenetur expedita sunt');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'dolor', 'eos qui et ipsum ipsam suscipit aut\nsed omnis non odio\nexpedita earum mollitia molestiae aut atque rem suscipit\nnam impedit esse');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'maxime', 'veritatis unde neque eligendi\nquae quod architecto quo neque vitae\nest illo sit tempora doloremque fugit quod\net et vel beatae sequi ullam sed tenetur perspiciatis');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'autem', 'enim et ex nulla\nomnis voluptas quia qui\nvoluptatem consequatur numquam aliquam sunt\ntotam recusandae id dignissimos aut sed asperiores deserunt');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'rem', 'ullam consequatur ut\nomnis quis sit vel consequuntur\nipsa eligendi ipsum molestiae et omnis error nostrum\nmolestiae illo tempore quia et distinctio');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'est', 'similique esse doloribus nihil accusamus\nomnis dolorem fuga consequuntur reprehenderit fugit recusandae temporibus\nperspiciatis cum ut laudantium\nomnis aut molestiae vel vero');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'quasi', 'eum sed dolores ipsam sint possimus debitis occaecati\ndebitis qui qui et\nut placeat enim earum aut odit facilis\nconsequatur suscipit necessitatibus rerum sed inventore temporibus consequatur');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'delectus', 'non et quaerat ex quae ad maiores\nmaiores recusandae totam aut blanditiis mollitia quas illo\nut voluptatibus voluptatem\nsimilique nostrum eum');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'iusto', 'odit magnam ut saepe sed non qui\ntempora atque nihil\naccusamus illum doloribus illo dolor\neligendi repudiandae odit magni similique sed cum maiores');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'aquo', 'alias dolor cumque\nimpedit blanditiis non eveniet odio maxime\nblanditiis amet eius quis tempora quia autem rem\na provident perspiciatis quia');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'ullam', 'debitis eius sed quibusdam non quis consectetur vitae\nimpedit ut qui consequatur sed aut in\nquidem sit nostrum et maiores adipisci atque\nquaerat voluptatem adipisci repudiandae');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'doloremque', 'deserunt eos nobis asperiores et hic\nest debitis repellat molestiae optio\nnihil ratione ut eos beatae quibusdam distinctio maiores\nearum voluptates et aut adipisci ea maiores voluptas maxime');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'qui', 'rerum ut et numquam laborum odit est sit\nid qui sint in\nquasi tenetur tempore aperiam et quaerat qui in\nrerum officiis sequi cumque quod');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'magnam', 'ea velit perferendis earum ut voluptatem voluptate itaque iusto\ntotam pariatur in\nnemo voluptatem voluptatem autem magni tempora minima in\nest distinctio qui assumenda accusamus dignissimos officia nesciunt nobis');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'id', 'nisi error delectus possimus ut eligendi vitae\nplaceat eos harum cupiditate facilis reprehenderit voluptatem beatae\nmodi ducimus quo illum voluptas eligendi\net nobis quia fugit');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'fuga', 'ad mollitia et omnis minus architecto odit\nvoluptas doloremque maxime aut non ipsa qui alias veniam\nblanditiis culpa aut quia nihil cumque facere et occaecati\nqui aspernatur quia eaque ut aperiam inventore');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'provident', 'debitis et eaque non officia sed nesciunt pariatur vel\nvoluptatem iste vero et ea\nnumquam aut expedita ipsum nulla in\nvoluptates omnis consequatur aut enim officiis in quam qui');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'explicabo', 'animi esse sit aut sit nesciunt assumenda eum voluptas\nquia voluptatibus provident quia necessitatibus ea\nrerum repudiandae quia voluptatem delectus fugit aut id quia\nratione optio eos iusto veniam iure');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eos', 'corporis rerum ducimus vel eum accusantium\nmaxime aspernatur a porro possimus iste omnis\nest in deleniti asperiores fuga aut\nvoluptas sapiente vel dolore minus voluptatem incidunt ex');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'enim', 'ut voluptatum aliquid illo tenetur nemo sequi quo facilis\nipsum rem optio mollitia quas\nvoluptatem eum voluptas qui\nunde omnis voluptatem iure quasi maxime voluptas nam');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'non', 'molestias id nostrum\nexcepturi molestiae dolore omnis repellendus quaerat saepe\nconsectetur iste quaerat tenetur asperiores accusamus ex ut\nnam quidem est ducimus sunt debitis saepe');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'commodi', 'odio fugit voluptatum ducimus earum autem est incidunt voluptatem\nodit reiciendis aliquam sunt sequi nulla dolorem\nnon facere repellendus voluptates quia\nratione harum vitae ut');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'eligendi', 'similique fugit est\nillum et dolorum harum et voluptate eaque quidem\nexercitationem quos nam commodi possimus cum odio nihil nulla\ndolorum exercitationem magnam ex et a et distinctio debitis');
